package api

import (
	"net/http"

	"github.com/N30A/korturl/internal/config"
	"github.com/N30A/korturl/internal/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type APIMux struct {
	router *http.ServeMux
	config config.Config
	pool   *pgxpool.Pool
}

func NewMux(config config.Config, pool *pgxpool.Pool) *APIMux {
	mux := &APIMux{
		router: http.NewServeMux(),
		config: config,
		pool:   pool,
	}
	mux.registerRoutes()
	return mux
}

func (m *APIMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}

func (m *APIMux) registerRoutes() {
	m.router.Handle("GET /health", middleware.ChainFunc(middleware.Logger)(m.healthHandler))
}
