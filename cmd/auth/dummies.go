package main

import (
	"context"
	"log/slog"
)

type DummyReader struct {
}

func (w *DummyReader) FetchUser(ctx context.Context, email string) (name, userID, password string, err error) {
	slog.Info("signing in user", "email", email, "service", "auth")

	return "name", "2", "$2a$10$3bAYWOIv0JgCj2xf9hf.beUioHU5jHIYED.hOxKLttWtNWFp7Aq/O", nil
}

type DummyWriter struct {
}

func (w *DummyWriter) StoreUser(ctx context.Context, name, email, password string) (userID string, err error) {
	slog.Info("registering user", "name", name, "email", email, "password", password, "service", "auth")

	return "2", nil
}
