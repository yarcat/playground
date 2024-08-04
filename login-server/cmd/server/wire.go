//go:build wireinject
// +build wireinject

package main

import (
	"login-server/pkg/crypto/argon2"
	ahttp "login-server/pkg/http"
	ajwt "login-server/pkg/jwt"
	"login-server/pkg/user"
	"login-server/pkg/user/inmem"

	"github.com/google/wire"
)

func InitializeServer() (*Server, error) {
	panic(wire.Build(
		wire.Bind(new(user.EmailReader), new(*inmem.Inmem)),
		wire.Bind(new(user.Hasher), new(user.Argon2Hasher)),
		wire.Bind(new(user.Factory), new(*inmem.Inmem)),
		wire.Value([]argon2.OptionFunc{}),
		wire.Struct(new(inmem.Inmem)),
		ahttp.NewLoginHandler,
		ahttp.NewRouter,
		ahttp.NewSignupHandler,
		ajwt.New,
		user.NewArgon2Hasher,
		NewHTTPLoginFromFlags,
		NewHTTPServerFromFlags,
		NewHTTPSignupFromFlags,
		NewSaltHasher,
		NewSecretProviderFromFlags,
		NewServer,
	))
}
