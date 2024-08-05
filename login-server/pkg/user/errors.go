package user

import (
	"errors"
	"fmt"
	t "login-server/types"
)

var (
	ErrNotFound = errors.New("not found")
)

// IsNotFound reports whether err is an ErrNotFound.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IDNotFound returns an error that reports that the ID is not found.
func IDNotFound(id t.UserID) error {
	return fmt.Errorf("id %q: %w", id, ErrNotFound)
}

// IDExists returns an error that reports that the ID already exists.
func IDExists(id t.UserID) error {
	return fmt.Errorf("id %q: already exists", id)
}

// EmailNotFound returns an error that reports that the email is not found.
func EmailNotFound(e t.Email) error {
	return fmt.Errorf("email %q: %w", e, ErrNotFound)
}

// EmailExists returns an error that reports that the email already exists.
func EmailExists(e t.Email) error {
	return fmt.Errorf("email %q: already exists", e)
}
