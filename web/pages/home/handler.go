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
	"github.com/N30A/korturl/web/pages"
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
	mux.Handle("GET /", middleware.Logger(http.HandlerFunc(h.Index)))
	mux.Handle("GET /{code}", middleware.Logger(http.HandlerFunc(h.Redirect)))
	mux.Handle("POST /", middleware.Logger(http.HandlerFunc(h.Shorten)))
	mux.Handle("POST /shorten-form", middleware.Logger(http.HandlerFunc(h.ShortenForm)))
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.PathValue("code"))
	if code == "" {
		http.Error(w, "short code is required", http.StatusBadRequest)
		return
	}

	if err := validation.URLCode(code); err != nil {
		render.Component(w, r, pages.Error("Ogiltig korturl", "Denna korturl är ogiltig."))
		return
	}

	if err := validation.ReservedURLCode(code); err != nil {
		render.Component(w, r, pages.Error("Resarverad korturl", "Denna korturl är resarverad."))
		return
	}

	redirectURL, err := h.urlService.GetRedirectURL(r.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			render.Component(w, r, pages.Error("Hittades ej", "Denna korturl finns inte."))
			return
		}
		slog.Error("failed to get redirect url", "error", err)
		render.Component(w, r, pages.Error("Något gick fel", "Ett oväntat fel uppstod, försök igen."))
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	render.Component(w, r, Page("korturl", ShortenFormData{}))
}

func (h *Handler) Shorten(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	alias := strings.TrimSpace(r.FormValue("alias"))

	redirectURL := strings.TrimSpace(r.FormValue("url"))
	if redirectURL == "" {
		render.Component(w, r, ShortenForm(ShortenFormData{Error: "URL:en får inte vara tom.", Alias: alias}))
		return
	}

	if !validation.RedirectURL(redirectURL) {
		render.Component(w, r, ShortenForm(ShortenFormData{Error: "Ogiltig HTTP eller HTTPS URL.", Alias: alias}))
		return
	}

	url, err := h.urlService.Create(r.Context(), redirectURL, alias)
	if err != nil {
		data := ShortenFormData{LongURL: redirectURL, Alias: alias}

		switch err {
		case service.ErrURLCodeAlreadyExists:
			data.AliasError = "Alias används redan."
		case validation.ErrURLCodeInvalid:
			data.AliasError = "Alias innehåller ogiltiga karaktärer, måste vara a-z, A-Z, 0-9, -"
		case validation.ErrURLCodeReserved:
			data.AliasError = "Alias is resarverat."
		case validation.ErrURLCodeTooShort:
			data.AliasError = fmt.Sprintf("Alias är för kort, måste minst vara %d characters.", validation.URLCodeMinLength)
		case validation.ErrURLCodeTooLong:
			data.AliasError = fmt.Sprintf("Alias är för kort, måste högst vara %d characters.", validation.URLCodeMaxLength)
		default:
			slog.Error("failed to create url", "error", err)
			data.Error = "Något gick fel, försök igen."
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
