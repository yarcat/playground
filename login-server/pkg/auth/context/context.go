package context

import (
	"context"
	"login-server/pkg/auth"
)

type key int

const authKey key = iota

func WithAuth(ctx context.Context, a *auth.Auth) context.Context {
	return context.WithValue(ctx, authKey, a)
}

func Auth(ctx context.Context) *auth.Auth {
	a, _ := ctx.Value(authKey).(*auth.Auth)
	return a
}
