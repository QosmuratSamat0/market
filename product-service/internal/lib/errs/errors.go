package errs

import "errors"

var (
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidInput  = errors.New("invalid input")

	ProductNotFound  = errors.New("product not found")
	CategoryNotFound = errors.New("category not found")
	Forbidden        = errors.New("forbidden")

	ErrUnauthorized = errors.New("unauthorized")
)
