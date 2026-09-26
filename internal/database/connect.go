package database

import (
	"context"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/N30A/korturl/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

func buildConnString(config config.DBConfig) string {
	params := config.Params
	if params != "" && !strings.HasPrefix(params, "?") {
		params = "?" + params
	}

	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s%s",
		url.QueryEscape(config.User),
		url.QueryEscape(config.Password),
		config.Host,
		config.Port,
		config.Name,
		params,
	)
}

func Connect(ctx context.Context, config config.DBConfig) (*pgxpool.Pool, error) {
	connString := buildConnString(config)

	pool, err := pgxpool.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		defer pool.Close()
		return nil, err
	}

	slog.Info("connected to database",
		"user", config.User,
		"name", config.Name,
		"host", config.Host,
		"port", config.Port,
		"params", config.Params,
	)
	return pool, nil
}
