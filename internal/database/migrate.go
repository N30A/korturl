package database

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/N30A/korturl/internal/config"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed migrations/*.sql
var fs embed.FS

func Migrate(config config.DBConfig) error {
	slog.Info("starting database migrations")

	source, err := iofs.New(fs, "migrations")
	if err != nil {
		return fmt.Errorf("migration: %w", err)
	}
	defer source.Close()

	connString := strings.Replace(buildConnString(config), "postgres://", "pgx5://", 1)

	m, err := migrate.NewWithSourceInstance("iofs", source, connString)
	if err != nil {
		return fmt.Errorf("migration: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("database migrations skipped", "reason", "no changes")
			return nil
		}
		return fmt.Errorf("migration: %w", err)
	}

	slog.Info("database migrations completed")
	return nil
}
