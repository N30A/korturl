package config

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL string
	DB      DBConfig
}

type DBConfig struct {
	Name     string
	User     string
	Password string
	Host     string
	Port     string
	Params   string
}

var requiredEnvs = []string{
	"DATABASE_NAME",
	"DATABASE_USER",
	"DATABASE_PASSWORD",
	"DATABASE_HOST",
	"DATABASE_PORT",
	"BASE_URL",
}

func getEnv(key string) string {
	return strings.TrimSpace(os.Getenv(key))
}

func requireEnvs() error {
	var builder strings.Builder
	for _, key := range requiredEnvs {
		if getEnv(key) == "" {
			fmt.Fprintf(&builder, "%s ", key)
		}
	}

	if builder.Len() > 0 {
		return fmt.Errorf("missing required environment variables: %s", builder.String())
	}

	return nil
}

func Load() (Config, error) {
	if err := godotenv.Load(); err != nil {
		slog.Info(".env not found, using environment instead")
	}

	if err := requireEnvs(); err != nil {
		return Config{}, err
	}

	return Config{
		BaseURL: getEnv("BASE_URL"),
		DB: DBConfig{
			Name:     getEnv("DATABASE_NAME"),
			User:     getEnv("DATABASE_USER"),
			Password: getEnv("DATABASE_PASSWORD"),
			Host:     getEnv("DATABASE_HOST"),
			Port:     getEnv("DATABASE_PORT"),
			Params:   getEnv("DATABASE_PARAMS"),
		},
	}, nil
}
