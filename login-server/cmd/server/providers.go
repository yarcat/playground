package main

import (
	"net/http"

	"login-server/pkg/crypto/argon2"
	ahttp "login-server/pkg/http"
	ajwt "login-server/pkg/jwt"
	"login-server/pkg/service"
	"login-server/pkg/user"
)

func NewHTTPServerFromFlags(h ahttp.Router) *http.Server {
	return &http.Server{
		Handler: h,
		Addr:    *addr,
	}
}

func NewHTTPLoginFromFlags(r user.EmailReader, sh user.SaltHasher) ahttp.Loginer {
	return service.NewLogin(r, sh, *issuer, *expiry)
}

func NewHTTPSignupFromFlags(r user.Factory, h user.Hasher) ahttp.Signuper {
	return service.NewSignup(r, h, *issuer, *expiry)
}

func NewSecretProviderFromFlags() (ajwt.SecretProvider, error) { return ajwt.NewArgon2(*password) }

type (
	// SaltHandler returns a new hasher given a salt.
	SaltHasherFunc func(user.SecretSalt) user.Hasher
)

// New implements the SaltHasher interface.
func (f SaltHasherFunc) New(s user.SecretSalt) user.Hasher { return f(s) }

func NewSaltHasher() user.SaltHasher {
	return SaltHasherFunc(func(ss user.SecretSalt) user.Hasher {
		return user.NewArgon2Hasher(argon2.WithSalt([]byte(ss)))
	})
}
