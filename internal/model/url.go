package model

import "time"

type URL struct {
	ID          string
	ShortID     string
	OriginalURL string
	UserID      string
	Visibility  string
	CreatedAt   time.Time
}
