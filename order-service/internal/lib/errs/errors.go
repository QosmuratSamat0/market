package errs

import "errors"

var (
	ErrInvalidInput = errors.New("invalid input")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")

	OrderNotFound   = errors.New("order not found")
	ProductNotFound = errors.New("product not found")
	OutOfStock      = errors.New("product out of stock")
)
