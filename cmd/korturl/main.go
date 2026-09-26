package main

import (
	"context"
	"log"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/N30A/korturl/api"
	"github.com/N30A/korturl/internal/config"
	"github.com/N30A/korturl/internal/database"
	"github.com/N30A/korturl/web"
)

const (
	host    = "127.0.0.1"
	port    = "8080"
	timeout = time.Second * 5
)

func main() {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	poolCtx, poolCancel := context.WithCancel(context.Background())
	pool, err := database.Connect(poolCtx, cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	webMux := web.NewMux(cfg, pool)
	apiMux := api.NewMux(cfg, pool)

	mux := http.NewServeMux()

	mux.Handle("/", webMux)
	mux.Handle("/api/", http.StripPrefix("/api", apiMux))

	server := &http.Server{
		Addr:    net.JoinHostPort(host, port),
		Handler: mux,
	}

	go func() {
		slog.Info("listening on", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Println(err)
			signals <- os.Interrupt
		}
	}()

	<-signals

	slog.Info("server shutting down", "timeout", timeout)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", "error", err)
	}

	pool.Close()
	poolCancel()
}
