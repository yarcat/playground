package http

import (
	"errors"
	"log/slog"
	"login-server/pkg/auth"
	"login-server/pkg/auth/context"
	ajwt "login-server/pkg/jwt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
)

func WithAuth(r *http.Request, a *auth.Auth) *http.Request {
	return r.WithContext(context.WithAuth(r.Context(), a))
}

func Auth(r *http.Request) *auth.Auth { return context.Auth(r.Context()) }

func OptionalAuth(j ajwt.JWT, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("auth_token") // TODO: Make configurable.
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				h.ServeHTTP(w, r)
				return
			}
		}
		t, err := j.Verify(token.Value)
		if err != nil {
			slog.WarnContext(r.Context(), "Failed to verify token", "error", err)
			h.ServeHTTP(w, r)
			return
		}
		if claims, ok := t.Claims.(jwt.MapClaims); ok {
			s := func(key string) string {
				v, _ := claims[key].(string)
				return v
			}
			r = WithAuth(r, auth.New(auth.Email(s("email")), auth.UserID(s("user_id"))))
		}
		h.ServeHTTP(w, r)
	}
}
