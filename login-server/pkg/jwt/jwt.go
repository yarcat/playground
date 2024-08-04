package jwt

import (
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type SecretProvider interface {
	Secret() []byte
}

type (
	// JWT provides a way to verify JWT tokens.
	JWT struct {
		SecretProvider
	}
)

// New returns a new JWT.
func New(secret SecretProvider) JWT {
	return JWT{secret}
}

// Verify parses and verifies the token. If the token is valid, it returns the token.
func (j JWT) Verify(token string) (Token, error) {
	t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("%w: %v", ErrSigningMethod, token.Header["alg"])
		}
		return j.Secret(), nil
	})
	return Token{t}, err
}
