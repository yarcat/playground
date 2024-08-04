package service

import (
	"login-server/pkg/jwt"
	"login-server/pkg/t"
	"login-server/pkg/user"
	"time"
)

type TokenOptions struct {
	Issuer string
	Expiry time.Duration
}

func newToken(u *user.User, opt TokenOptions) (jwt.Token, error) {
	type Claims struct {
		t.Email   `json:"email"`
		user.ID   `json:"id"`
		user.Role `json:"role"`
		jwt.RegisteredClaims
	}
	now := time.Now()
	return jwt.NewTokenWithClaims(Claims{
		ID:    u.ID,
		Email: u.Email,
		Role:  u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    opt.Issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(opt.Expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}), nil
}
