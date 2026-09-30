package home

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/N30A/korturl/internal/config"
	"github.com/N30A/korturl/internal/middleware"
	"github.com/N30A/korturl/internal/service"
	"github.com/N30A/korturl/internal/validation"
	"github.com/jackc/pgx/v5"
)

type Handler struct {
	baseURL    string
	urlService *service.URLService
}

func NewHandler(config config.Config, urlService *service.URLService) *Handler {
	return &Handler{
		baseURL:    config.BaseURL,
		urlService: urlService,
	}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /", middleware.Logger(http.HandlerFunc(h.Index)))
	mux.Handle("GET /{code}", middleware.Logger(http.HandlerFunc(h.Redirect)))
	mux.Handle("POST /shorten", middleware.Logger(http.HandlerFunc(h.Shorten)))
	mux.Handle("POST /shorten/form", middleware.Logger(http.HandlerFunc(h.ShortenForm)))
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.PathValue("code"))
	if code == "" {
		http.Error(w, "short code is required", http.StatusBadRequest)
		return
	}

	if err := validation.URLCode(code); err != nil {
		http.Error(w, "invalid short code", http.StatusBadRequest)
		return
	}

	redirectURL, err := h.urlService.GetRedirectURL(r.Context(), code)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "url does not exist", http.StatusNotFound)
			return
		}
		slog.Error("failed to get redirect url", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	err := Page("korturl", ShortenFormData{}).Render(r.Context(), w)
	if err != nil {
		slog.Error("failed to render home page", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	redirectURL := strings.TrimSpace(r.FormValue("url"))
	if redirectURL == "" {
		data := ShortenFormData{Error: "URL must not be empty."}
		if err := ShortenForm(data).Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	if !validation.RedirectURL(redirectURL) {
		data := ShortenFormData{Error: "URL must be a valid HTTP or HTTPS URL."}
		if err := ShortenForm(data).Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	alias := strings.TrimSpace(r.FormValue("alias"))

	url, err := h.urlService.Create(r.Context(), redirectURL, alias)
	if err != nil {
		data := ShortenFormData{LongURL: redirectURL}

		switch err {
		case service.ErrURLCodeAlreadyExists:
			data.AliasError = "Alias is already in use"
		case validation.ErrURLCodeInvalid:
			data.AliasError = "Alias contains invalid characters, must be a-z, A-Z, 0-9, -"
		case validation.ErrURLCodeReserved:
			data.AliasError = "Alias is reserved"
		case validation.ErrURLCodeTooShort:
			data.AliasError = fmt.Sprintf("Alias is too short, must be at least %d characters", validation.URLCodeMinLength)
		case validation.ErrURLCodeTooLong:
			data.AliasError = fmt.Sprintf("Alias is too long, must be at most %d characters", validation.URLCodeMaxLength)
		default:
			slog.Error("failed to create url", "error", err)
			data.Error = "Something went wrong, please try again."
		}

		if err := ShortenForm(data).Render(r.Context(), w); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}

	shortURL := h.baseURL + "/" + url.Code

	if err := ShortenResult(redirectURL, shortURL).Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) ShortenForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data := ShortenFormData{LongURL: r.FormValue("url")}
	if err := ShortenForm(data).Render(r.Context(), w); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
