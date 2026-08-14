package model

import "time"

type User struct {
	ID            string
	GoogleSubject string
	Email         string
	Name          string
	AvatarURL     *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
