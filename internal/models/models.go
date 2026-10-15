package models

import "time"

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
