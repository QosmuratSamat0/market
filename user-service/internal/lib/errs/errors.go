package errs

import "errors"

var (
	// General
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")
	ErrInvalidInput  = errors.New("invalid input")

	// User
	UserNotFound       = errors.New("user not found")
	InvalidUserID      = errors.New("invalid user id")
	EmailAlreadyExists = errors.New("email already exists")
	Forbidden          = errors.New("forbidden")

	// Auth
	ErrUnauthorized = errors.New("unauthorized")
)
