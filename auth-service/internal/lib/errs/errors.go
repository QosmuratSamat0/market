package errs

import "errors"

var (
	// General
	ErrNotFound      = errors.New("not found")
	ErrAlreadyExists = errors.New("already exists")

	// User
	UserNotFound          = errors.New("user not found")
	InvalidUserID         = errors.New("invalid user id")
	InvalidUserEmail      = errors.New("invalid user email")
	InvalidRole           = errors.New("invalid user role")
	EmailAlreadyExists    = errors.New("email already exists")
	Forbidden             = errors.New("forbidden")
	CannotDeleteSelf      = errors.New("cannot delete self")
	CannotDeleteLastAdmin = errors.New("cannot delete last admin")
	CannotDemoteLastAdmin = errors.New("cannot demote last admin")

	// Refresh Token
	RefreshTokenNotFound   = errors.New("refresh token not found")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")

	// Auth
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidInput        = errors.New("invalid input")
	ErrInvalidCredentials  = errors.New("invalid credentials")
	ErrTooManyLoginAttempt = errors.New("too many login attempts, please try again later")
)
