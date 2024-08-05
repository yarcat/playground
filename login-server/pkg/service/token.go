package service

import (
	"login-server/pkg/jwt"
	"login-server/pkg/user"
	t "login-server/types"
	"time"
)

type TokenOptions struct {
	Issuer string
	Expiry time.Duration
}

func newToken(u *user.User, opt TokenOptions) (jwt.Token, error) {
	type Claims struct {
		t.Email  `json:"email"`
		t.UserID `json:"id"`
		t.Role   `json:"role"`
		jwt.RegisteredClaims
	}
	now := time.Now()
	return jwt.NewTokenWithClaims(Claims{
		UserID: u.UserID,
		Email:  u.Email,
		Role:   u.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    opt.Issuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(opt.Expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}), nil
}
