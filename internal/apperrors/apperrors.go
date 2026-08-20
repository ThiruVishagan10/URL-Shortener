package errors

import "errors"

var (
	ErrURLNotFound     = errors.New("URL not found")
	ErrDuplicateURL    = errors.New("URL already exists")
	ErrUserNotFound    = errors.New("User not found")
	ErrDuplicateUser   = errors.New("user already exists")
	ErrSessionNotFound = errors.New("session expired Found")
)
