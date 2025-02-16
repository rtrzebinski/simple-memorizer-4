package main

import (
	"context"
	"log/slog"
)

type DummyWriter struct {
}

func (w *DummyWriter) Register(ctx context.Context, name, email, password string) (userID string, err error) {
	slog.Info("registering user", "name", name, "email", email)

	return "userID", nil
}
