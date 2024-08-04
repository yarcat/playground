package jwt

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type (
	// Token represents a JWT token.
	Token struct {
		*jwt.Token
	}
	RegisteredClaims = jwt.RegisteredClaims
	// Claims represent any form of a JWT Claims Set according to
	// https://datatracker.ietf.org/doc/html/rfc7519#section-4. In order to have a
	// common basis for validation, it is required that an implementation is able to
	// supply at least the claim names provided in
	// https://datatracker.ietf.org/doc/html/rfc7519#section-4.1 namely `exp`,
	// `iat`, `nbf`, `iss`, `sub` and `aud`.
	Claims = jwt.Claims
	// NumericDate represents a JSON numeric date value, as referenced at
	// https://datatracker.ietf.org/doc/html/rfc7519#section-2.
	NumericDate = jwt.NumericDate
)

// NewTokenWithClaims returns a new Token with the given claims.
func NewTokenWithClaims(claims jwt.Claims) Token {
	return Token{jwt.NewWithClaims(jwt.SigningMethodHS256, claims)}
}

// NewNumericDate constructs a new *NumericDate from a standard library time.Time struct.
func NewNumericDate(t time.Time) *NumericDate { return jwt.NewNumericDate(t) }
