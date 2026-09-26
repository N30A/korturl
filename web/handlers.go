package web

import (
	"errors"
	"log/slog"
	"net/http"
	"strings"

	"github.com/N30A/korturl/internal/httpjson"
	"github.com/N30A/korturl/internal/shortener"
	"github.com/jackc/pgx/v5"
)

func (m *WebMux) redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimSpace(r.PathValue("code"))
	if code == "" {
		http.Error(w, "short code is required", http.StatusBadRequest)
		return
	}

	if !shortener.ValidateShortCode(code) {
		http.Error(w, "invalid short code", http.StatusBadRequest)
		return
	}

	var redirectURL string
	if err := m.pool.QueryRow(r.Context(), "SELECT redirect_url FROM urls WHERE short_code = $1", code).Scan(&redirectURL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "url does not exist", http.StatusNotFound)
			return
		}
		slog.Error("failed to scan redirect url", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

type shortenURLRequest struct {
	URL string `json:"url"`
}

type shortenURLResponse struct {
	ShortCode   string `json:"short_code"`
	ShortURL    string `json:"short_url"`
	RedirectURL string `json:"redirect_url"`
}

func (m *WebMux) shortenURLHandler(w http.ResponseWriter, r *http.Request) {
	request, err := httpjson.Decode[shortenURLRequest](r)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	shortCode, err := shortener.GenerateUniqueShortCode(r.Context(), m.pool)
	if err != nil {
		slog.Error("failed to generate unique short code", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	query := "INSERT INTO urls (short_code, redirect_url) VALUES ($1, $2)"

	if _, err := m.pool.Exec(r.Context(), query, shortCode, request.URL); err != nil {
		slog.Error("failed to insert url", "error", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	shortURL := m.config.BaseURL + "/" + shortCode

	httpjson.Write(w, http.StatusCreated, shortenURLResponse{
		ShortCode:   shortCode,
		ShortURL:    shortURL,
		RedirectURL: request.URL,
	})
}
