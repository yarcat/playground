package main

import (
	"context"
	"errors"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	addr     = flag.String("addr", ":8080", "server address")
	password = flag.String("password", "password", "password to use for JWT")
	issuer   = flag.String("issuer", "login-server", "JWT issuer")
	expiry   = flag.Duration("expiry", 24*time.Hour, "JWT expiry duration")
)

func main() {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	s, err := InitializeServer()
	if err != nil {
		panic(err)
	}

	slog.Info("Server started", "addr", *addr)
	if err := Run(ctx, s); err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}
	slog.Info("Server stopped")
}
