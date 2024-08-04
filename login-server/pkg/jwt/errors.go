package jwt

import "errors"

var (
	// ErrSigningMethod is returned when the signing method is unexpected.
	ErrSigningMethod = errors.New("invalid signing method")
)
