package service

import (
	"context"
	"fmt"
	"login-server/pkg/jwt"
	"login-server/pkg/t"
	"login-server/pkg/user"
	"time"
)

// Signup is the signup service, it creates new users.
type Signup struct {
	user.Factory
	user.Hasher
	TokenOptions
}

// NewSignup returns a new Signup service.
func NewSignup(f user.Factory, h user.Hasher, iss string, exp time.Duration) *Signup {
	return &Signup{
		Factory: f,
		Hasher:  h,
		TokenOptions: TokenOptions{
			Issuer: iss,
			Expiry: exp,
		},
	}
}

// Signup creates a new user.
func (s *Signup) Signup(ctx context.Context, n user.Name, e t.Email, p t.Password) (jwt.Token, error) {
	secret, salt, err := s.Hash(p)
	if err != nil {
		return jwt.Token{}, fmt.Errorf("signup: %w", err)
	}
	u, err := s.New(ctx, n, e, secret, salt)
	if err != nil {
		return jwt.Token{}, fmt.Errorf("signup: %w", err)
	}
	return newToken(u, s.TokenOptions)
}
