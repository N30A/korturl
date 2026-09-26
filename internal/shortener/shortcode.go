package shortener

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	alphabet = "abcdefghjkmnpqrstuvwxyz23456789"
	length   = 8
)

func GenerateUniqueShortCode(ctx context.Context, pool *pgxpool.Pool) (string, error) {
	var (
		lastErr    error
		maxRetries int = 3
	)

	for range maxRetries {
		shortCode, err := randomShortCode(length)
		if err != nil {
			lastErr = err
			continue
		}

		exists, err := ShortCodeExists(ctx, pool, shortCode)
		if err != nil {
			lastErr = err
			continue
		}

		if exists {
			continue
		}

		return shortCode, nil
	}

	return "", fmt.Errorf("failed to generate unique short code after %d retries: %w", maxRetries, lastErr)
}

func ShortCodeExists(ctx context.Context, pool *pgxpool.Pool, shortCode string) (bool, error) {
	query := "SELECT EXISTS (SELECT 1 FROM urls where short_code = $1)"

	var exists bool
	if err := pool.QueryRow(ctx, query, shortCode).Scan(&exists); err != nil {
		return false, err
	}

	return exists, nil
}

func ValidateShortCode(code string) bool {
	for _, c := range code {
		if !strings.ContainsRune(alphabet, c) {
			return false
		}
	}

	return true
}

func randomShortCode(length int) (string, error) {
	result := make([]byte, length)

	for i := range result {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(alphabet))))
		if err != nil {
			return "", err
		}

		result[i] = alphabet[n.Int64()]
	}

	return string(result), nil
}
