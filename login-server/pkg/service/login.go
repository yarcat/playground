package service

import (
	"context"
	"errors"
	"fmt"
	"login-server/pkg/jwt"
	"login-server/pkg/t"
	"login-server/pkg/user"
	"time"
)

type (
	// Login logs the user in using the email and password.
	Login struct {
		user.EmailReader
		user.SaltHasher
		TokenOptions
	}
)

// NewLogin returns a new login service.
func NewLogin(r user.EmailReader, sh user.SaltHasher, iss string, exp time.Duration) *Login {
	return &Login{
		EmailReader: r,
		SaltHasher:  sh,
		TokenOptions: TokenOptions{
			Issuer: iss,
			Expiry: exp,
		},
	}
}

// Login logs the user in using the email and password.
func (l *Login) Login(ctx context.Context, e t.Email, p t.Password) (jwt.Token, error) {
	u, err := l.FromEmail(ctx, e)
	if err != nil {
		return jwt.Token{}, fmt.Errorf("login: %w", err)
	}
	hash, _, err := l.New(u.SecretSalt).Hash(p)
	if err != nil {
		return jwt.Token{}, fmt.Errorf("login: %w", err)
	} else if hash != u.SecretHash {
		return jwt.Token{}, errors.New("login: wrong password") // TODO: Use a typed error.
	}
	return newToken(u, l.TokenOptions)
}
