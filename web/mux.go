package web

import (
	"net/http"

	"github.com/N30A/korturl/internal/config"
	"github.com/N30A/korturl/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WebMux struct {
	router      *http.ServeMux
	config      config.Config
	pool        *pgxpool.Pool
	rateLimiter *middleware.RateLimiter
}

func NewMux(config config.Config, pool *pgxpool.Pool) *WebMux {
	mux := &WebMux{
		router:      http.NewServeMux(),
		config:      config,
		pool:        pool,
		rateLimiter: middleware.NewRateLimiter(middleware.DefaultRateLimitConfig),
	}
	mux.registerRoutes()
	return mux
}

func (m *WebMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}

func (m *WebMux) registerRoutes() {
	m.router.Handle("GET /{code}", middleware.ChainFunc(
		middleware.Logger,
		middleware.RateLimit(m.rateLimiter),
	)(m.redirectHandler))
	m.router.Handle("POST /shorten-url", middleware.ChainFunc(
		middleware.Logger,
		middleware.RateLimit(m.rateLimiter),
	)(m.shortenURLHandler))
}
