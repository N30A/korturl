package models

import "time"

type UserRole string

const (
	RoleUser  UserRole = "user"
	RoleAdmin UserRole = "admin"
)

type User struct {
	ID           int
	Username     string
	PasswordHash string
	Role         UserRole
}

type APIKey struct {
	ID        int
	UserID    int
	Name      string
	KeyHash   string
	CreatedAt time.Time
	ExpiresAt *time.Time
	RevokedAt *time.Time
}

type URL struct {
	ID          int
	Code        string
	RedirectURL string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Visit struct {
	ID        int64
	URLID     int
	VisitedAt time.Time
}
