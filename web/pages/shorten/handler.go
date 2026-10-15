package shorten

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/N30A/korturl/internal/config"
	"github.com/N30A/korturl/internal/middleware"
	"github.com/N30A/korturl/internal/service"
	"github.com/N30A/korturl/internal/validation"
	"github.com/N30A/korturl/web/render"
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
	mux.Handle("GET /shorten", middleware.Logger(http.HandlerFunc(h.Index)))
	mux.Handle("POST /shorten", middleware.Logger(http.HandlerFunc(h.Shorten)))
	mux.Handle("POST /shorten/form", middleware.Logger(http.HandlerFunc(h.ShortenForm)))
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	render.Component(w, r, Page("korturl", ShortenFormData{}))
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	redirectURL := strings.TrimSpace(r.FormValue("url"))
	if redirectURL == "" {
		render.Component(w, r, ShortenForm(ShortenFormData{Error: "URL must not be empty."}))
		return
	}

	if !validation.RedirectURL(redirectURL) {
		render.Component(w, r, ShortenForm(ShortenFormData{Error: "URL must be a valid HTTP or HTTPS URL."}))
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

		render.Component(w, r, ShortenForm(data))
		return
	}

	shortURL := h.baseURL + "/" + url.Code

	render.Component(w, r, ShortenResult(redirectURL, shortURL))
}

func (h *Handler) ShortenForm(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	render.Component(w, r, ShortenForm(ShortenFormData{LongURL: r.FormValue("url")}))
}
