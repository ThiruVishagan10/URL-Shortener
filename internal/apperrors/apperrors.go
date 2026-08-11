package errors

import "errors"

var (
	ErrURLNotFound = errors.New("URL not found")
	ErrDuplicateURL = errors.New("URL already exists")
)