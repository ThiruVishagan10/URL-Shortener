package model

import "time"

type URL struct {
	ID          string
	ShortID     string
	OriginalURL string
	CreatedAt   time.Time
}
