package user

import (
	"errors"
	"fmt"
	"login-server/pkg/t"
)

var (
	ErrNotFound = errors.New("not found")
)

// IsNotFound reports whether err is an ErrNotFound.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IDNotFound returns an error that reports that the ID is not found.
func IDNotFound(id ID) error {
	return fmt.Errorf("id %q: %w", id, ErrNotFound)
}

// EmailNotFound returns an error that reports that the email is not found.
func EmailNotFound(e t.Email) error {
	return fmt.Errorf("email %q: %w", e, ErrNotFound)
}
