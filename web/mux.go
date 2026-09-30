package web

import (
	"net/http"

	"github.com/N30A/korturl/internal/config"
	"github.com/N30A/korturl/internal/service"
	"github.com/N30A/korturl/web/home"
	"github.com/N30A/korturl/web/static"
)

type WebMux struct {
	router     *http.ServeMux
	config     config.Config
	urlService *service.URLService

	homeHandler *home.Handler
}

func NewMux(config config.Config, urlService *service.URLService) *WebMux {
	mux := &WebMux{
		router:     http.NewServeMux(),
		config:     config,
		urlService: urlService,

		homeHandler: home.NewHandler(config, urlService),
	}
	mux.registerRoutes()
	return mux
}

func (m *WebMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.router.ServeHTTP(w, r)
}

func (m *WebMux) registerRoutes() {
	m.router.Handle("GET /static/", http.StripPrefix("/static", http.FileServer(http.FS(static.Files))))
	m.homeHandler.RegisterRoutes(m.router)
	// m.router.Handle("GET /setup", protected(m.setupHandler))
	// m.router.Handle("POST /setup", protected(m.setupHandler))
}
