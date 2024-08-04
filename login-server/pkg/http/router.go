package http

import (
	ajwt "login-server/pkg/jwt"
	"net/http"
)

type Router struct {
	*http.ServeMux
}

func NewRouter(
	j ajwt.JWT,
	login *LoginHandler,
	signup *SignupHandler,
) Router {
	s := http.NewServeMux()
	s.HandleFunc("/", OptionalAuth(j, Template("index.html")))
	s.HandleFunc("/login", OptionalAuth(j, Template("login.html")))
	s.Handle("POST /login", login)
	s.Handle("POST /signup", signup)
	return Router{s}
}
