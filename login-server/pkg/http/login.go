package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	ajwt "login-server/pkg/jwt"
	t "login-server/types"
)

type (
	// Loginer logs the user in using the email and password.
	Loginer interface {
		Login(context.Context, t.Email, t.Password) (ajwt.Token, error)
	}

	// LoginHandler handles the login requests.
	LoginHandler struct {
		Loginer
		ajwt.SecretProvider
	}
)

// NewLoginHandler returns a new login handler.
func NewLoginHandler(l Loginer, p ajwt.SecretProvider) *LoginHandler {
	return &LoginHandler{Loginer: l, SecretProvider: p}
}

// ServeHTTP handles the login request.
func (lh *LoginHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, password := r.PostFormValue("email"), r.PostFormValue("password")
	if email == "" || password == "" {
		http.Error(w, "empty email or password", http.StatusBadRequest)
		return
	}

	token, err := lh.Login(r.Context(), t.Email(email), t.Password(password))
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to login", "error", err)
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	if setAuthTokenCookie(w, token, lh.Secret()) != nil {
		return
	}

	// TODO: Support meaningfull redirects.
	http.Redirect(w, r, r.URL.Path, http.StatusFound) // Redirect to the same page.
}

func setAuthTokenCookie(w http.ResponseWriter, token ajwt.Token, secret []byte) error {
	signed, err := token.SignedString(secret)
	if err != nil {
		slog.Error("Failed to sign the token", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "auth_token",
		Value:    signed,
		Expires:  time.Now().Add(24 * time.Hour),
		Path:     "/",
		Domain:   "localhost",
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
	return nil
}

type (
	// Signuper signs up a new user.
	Signuper interface {
		Signup(context.Context, t.Name, t.Email, t.Password) (ajwt.Token, error)
	}
	// SignupHandler handles the signup requests.
	SignupHandler struct {
		Signuper
		ajwt.SecretProvider
	}
)

// NewSignupHandler returns a new signup handler.
func NewSignupHandler(s Signuper, p ajwt.SecretProvider) *SignupHandler {
	return &SignupHandler{Signuper: s, SecretProvider: p}
}

// ServeHTTP handles the signup request.
func (sh *SignupHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	email, password := r.PostFormValue("email"), r.PostFormValue("password")
	if email == "" || password == "" {
		http.Error(w, "empty email or password", http.StatusBadRequest)
		return
	}

	u, err := sh.Signup(r.Context(), "" /*name*/, t.Email(email), t.Password(password))
	if err != nil {
		slog.ErrorContext(r.Context(), "Failed to create a new user", "error", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if setAuthTokenCookie(w, u, sh.Secret()) != nil {
		return
	}

	// TODO: Support meaningfull redirects.
	http.Redirect(w, r, "/", http.StatusFound)

}
