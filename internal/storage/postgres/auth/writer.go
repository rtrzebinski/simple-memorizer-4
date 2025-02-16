package auth

import (
	"context"
	"database/sql"
)

type Writer struct {
	db *sql.DB
}

func NewWriter(db *sql.DB) *Writer {
	return &Writer{db: db}
}

func (w *Writer) StoreUser(ctx context.Context, name, email, password string) (userID string, err error) {
	const query = `INSERT INTO "user" (name, email, password) VALUES ($1, $2, $3) RETURNING id;`

	row := w.db.QueryRowContext(ctx, query, name, email, password)

	err = row.Scan(&userID)
	if err != nil {
		return "", err
	}

	return userID, nil
}
