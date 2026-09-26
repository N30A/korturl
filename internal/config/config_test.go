package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	t.Setenv("DATABASE_NAME", "testdb")

	got := getEnv("DATABASE_NAME")
	want := "testdb"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGetEnvNoValue(t *testing.T) {
	t.Setenv("DATABASE_NAME", "")

	got := getEnv("DATABASE_NAME")
	want := ""

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGetEnvNoVariableSet(t *testing.T) {
	os.Unsetenv("DATABASE_NAME")

	got := getEnv("DATABASE_NAME")
	want := ""

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestGetEnvWithSpace(t *testing.T) {
	t.Setenv("DATABASE_NAME", " testdb ")

	got := getEnv("DATABASE_NAME")
	want := "testdb"

	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestRequireEnvsAllPresent(t *testing.T) {
	for _, key := range requiredEnvs {
		t.Setenv(key, "testvalue")
	}

	err := requireEnvs()
	if err != nil {
		t.Errorf("Required env variables not set: %v", err)
	}
}

func TestRequireEnvsSomeMissing(t *testing.T) {
	t.Setenv("DATABASE_NAME", "testdb")
	t.Setenv("DATABASE_USER", "testuser")
	t.Setenv("DATABASE_PASSWORD", "testpass")

	err := requireEnvs()
	if err == nil {
		t.Errorf("Expected error for missing env variables")
	}
}

func TestLoad(t *testing.T) {
	t.Setenv("DATABASE_NAME", "testdb")
	t.Setenv("DATABASE_USER", "testuser")
	t.Setenv("DATABASE_PASSWORD", "testpass")
	t.Setenv("DATABASE_HOST", "localhost")
	t.Setenv("DATABASE_PORT", "5432")
	t.Setenv("BASE_URL", "http://localhost:8080")

	_, err := Load()
	if err != nil {
		t.Errorf("Failed to load config: %v", err)
	}
}

func TestLoadWithEnvFile(t *testing.T) {
	for _, key := range []string{"DATABASE_NAME", "DATABASE_USER", "DATABASE_PASSWORD", "DATABASE_HOST", "DATABASE_PORT"} {
		os.Unsetenv(key)
	}

	content := `
		DATABASE_NAME=testdb
		DATABASE_USER=testuser
		DATABASE_PASSWORD=testpass
		DATABASE_HOST=localhost
		DATABASE_PORT=5432
		BASE_URL=http://localhost:8080
	`

	filename := ".env"

	file, err := os.Create(filename)
	if err != nil {
		t.Fatalf("Failed to create .env file: %v", err)
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		t.Fatalf("Failed to write to .env file: %v", err)
	}
	defer os.Remove(filename)

	_, err = Load()
	if err != nil {
		t.Errorf("Failed to load config: %v", err)
	}
}
