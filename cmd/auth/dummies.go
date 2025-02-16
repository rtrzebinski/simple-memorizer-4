package main

import (
	"context"
	"log/slog"
)

type DummyReader struct {
}

func (w *DummyReader) SignIn(ctx context.Context, email, password string) (name, userID string, err error) {
	slog.Info("signing in user", "email", email)

	return "name", "userID", nil
}

type DummyWriter struct {
}

func (w *DummyWriter) Register(ctx context.Context, name, email, password string) (userID string, err error) {
	slog.Info("registering user", "name", name, "email", email)

	return "userID", nil
}
