package postgres

import (
	"context"
	"database/sql"
)

type AuthWriter struct {
	db *sql.DB
}

func NewAuthWriter(db *sql.DB) *AuthWriter {
	return &AuthWriter{db: db}
}

func (w *AuthWriter) StoreUser(ctx context.Context, name, email, password string) (userID string, err error) {
	const query = `INSERT INTO "user" (name, email, password) VALUES ($1, $2, $3) RETURNING id;`

	row := w.db.QueryRowContext(ctx, query, name, email, password)

	err = row.Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
