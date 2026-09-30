package service

import (
	"context"
	"crypto/rand"
	"errors"
	"math/big"
	"strings"

	"github.com/N30A/korturl/internal/database"
	"github.com/N30A/korturl/internal/models"
	"github.com/N30A/korturl/internal/validation"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	codeAlphabet = "abcdefghjkmnpqrstuvwxyz23456789"
	codeLength   = 6
)

var (
	ErrURLCodeAlreadyExists     = errors.New("url with code already exists")
	ErrURLCodeReachedMaxRetries = errors.New("url code reached max retries")
	ErrURLNotFound              = errors.New("url was not found")
)

type URLService struct {
	pool *pgxpool.Pool
}

func NewURLService(pool *pgxpool.Pool) *URLService {
	return &URLService{pool: pool}
}

// Create creates a new URL. If an alias is provided, it is used as the URL code;
// otherwise, a random code is generated. Random code generation is retried up to
// a maximum number of attempts if a generated code already exists.
//
// Specific errors:
//
//	ErrURLCodeAlreadyExists: if the provided alias already exists
//	ErrURLCodeTooShort: if the code is too short
//	ErrURLCodeTooLong: if the code is too long
//	ErrURLCodeInvalid: if the code contains invalid characters
//	ErrURLCodeReachedMaxRetries: if a unique random code could not be generated within the maximum number of retries
func (s *URLService) Create(ctx context.Context, redirectURL, alias string) (models.URL, error) {
	redirectURL = strings.TrimSpace(redirectURL)
	alias = strings.TrimSpace(alias)

	if alias != "" {
		if err := validation.URLCode(alias); err != nil {
			return models.URL{}, err
		}

		if err := validation.ReservedURLCode(alias); err != nil {
			return models.URL{}, err
		}

		return s.insertURL(ctx, alias, redirectURL)
	}

	const maxRetries = 10

	for range maxRetries {
		code, err := generateRandomCode(codeLength)
		if err != nil {
			return models.URL{}, err
		}

		if err := validation.URLCode(code); err != nil {
			return models.URL{}, err
		}

		url, err := s.insertURL(ctx, code, redirectURL)
		if err == nil {
			return url, nil
		}

		if !errors.Is(err, ErrURLCodeAlreadyExists) {
			return models.URL{}, err
		}
	}

	return models.URL{}, ErrURLCodeReachedMaxRetries
}

// insertURL inserts the URL with its code to the database.
//
// It returns the inserted URL, or ErrURLCodeAlreadyExists if a URL with the given code already exists.
// Any other errors encountered while querying the database is returned unchanged.
func (s *URLService) insertURL(ctx context.Context, code, redirectURL string) (models.URL, error) {
	query := `
		INSERT INTO urls (code, redirect_url)
		VALUES ($1, $2)
		RETURNING id, code, redirect_url, created_at, updated_at
	`

	var model models.URL

	err := s.pool.QueryRow(ctx, query, code, redirectURL).Scan(
		&model.ID, &model.Code, &model.RedirectURL,
		&model.CreatedAt, &model.UpdatedAt,
	)
	if err != nil {
		if database.IsError(err, database.ErrUniqueViolation) {
			return models.URL{}, ErrURLCodeAlreadyExists
		}
		return models.URL{}, err
	}

	return model, nil
}

func (s *URLService) GetRedirectURL(ctx context.Context, code string) (string, error) {
	query := "SELECT redirect_url FROM urls WHERE code = $1"

	var redirectURL string
	if err := s.pool.QueryRow(ctx, query, code).Scan(&redirectURL); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrURLNotFound
		}
		return "", err
	}

	return redirectURL, nil
}

// generateRandomCode generates a cryptographically secure random code with the
// specified length. Each character is randomly selected from codeAlphabet.
//
// The generated code is not guaranteed to be unique, so callers must handle
// potential collisions when uniqueness is required.
//
// It returns the generated code, or an error if the random number generator
// fails. Such an error is unexpected and typically indicates a problem with
// the system's cryptographic random number source.
func generateRandomCode(length int) (string, error) {
	code := make([]byte, length)
	for i := range code {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(codeAlphabet))))
		if err != nil {
			return "", err
		}

		code[i] = codeAlphabet[n.Int64()]
	}

	return string(code), nil
}
