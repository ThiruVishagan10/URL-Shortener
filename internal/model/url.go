package model

import "time"

type URL struct {
	ID          string
	ShortID     string
	OriginalURL string
	UserID      string
	CreatedAt   time.Time
}
