package database

import (
	"testing"

	"github.com/N30A/korturl/internal/config"
)

func TestBuildConnString(t *testing.T) {
	config := config.DBConfig{
		User:     "testuser",
		Password: "testpass",
		Host:     "localhost",
		Port:     "5432",
		Name:     "testdb",
		Params:   "sslmode=disable",
	}

	got := buildConnString(config)
	want := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"

	if got != want {
		t.Errorf("Expected connection string '%s', got '%s'", want, got)
	}
}

func TestBuildConnStringWithoutParams(t *testing.T) {
	config := config.DBConfig{
		User:     "testuser",
		Password: "testpass",
		Host:     "localhost",
		Port:     "5432",
		Name:     "testdb",
	}

	got := buildConnString(config)
	want := "postgres://testuser:testpass@localhost:5432/testdb"

	if got != want {
		t.Errorf("Expected connection string '%s', got '%s'", want, got)
	}
}

func TestBuildConnStringWithParamsHavingPrefix(t *testing.T) {
	config := config.DBConfig{
		User:     "testuser",
		Password: "testpass",
		Host:     "localhost",
		Port:     "5432",
		Name:     "testdb",
		Params:   "?sslmode=disable",
	}

	got := buildConnString(config)
	want := "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable"

	if got != want {
		t.Errorf("Expected connection string '%s', got '%s'", want, got)
	}
}
