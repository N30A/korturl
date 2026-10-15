package home

import (
	"errors"
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
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.PathValue("code"))
	if code == "" {
		http.Error(w, "short code is required", http.StatusBadRequest)
		return
	}

	if err := validation.URLCode(code); err != nil {
		render.Component(w, r, pages.Error("Invalid short URL", "This short URL is not valid."))
		return
	}

	if err := validation.ReservedURLCode(code); err != nil {
		render.Component(w, r, pages.Error("Reserved short URL", "This short URL cannot be used."))
		return
	}

	redirectURL, err := h.urlService.GetRedirectURL(r.Context(), code)
	if err != nil {
		if errors.Is(err, service.ErrURLNotFound) {
			render.Component(w, r, pages.Error("Not Found", "This short URL does not exist."))
			return
		}
		slog.Error("failed to get redirect url", "error", err)
		render.Component(w, r, pages.Error("Something went wrong", "An unexpected error occurred, please try again later."))
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

func (h *Handler) Index(w http.ResponseWriter, r *http.Request) {
	render.Component(w, r, Page("korturl"))
}
