package postgres

import (
	"context"
	"database/sql"
)

type AuthReader struct {
	db *sql.DB
}

func NewAuthReader(db *sql.DB) *AuthReader {
	return &AuthReader{db: db}
}

func (r *AuthReader) FetchUser(ctx context.Context, email string) (name, userID, password string, err error) {
	const query = `SELECT name, id, password FROM "user" WHERE email = $1;`

	row := r.db.QueryRowContext(ctx, query, email)

	err = row.Scan(&name, &userID, &password)
	if err != nil {
		return "", "", "", err
	}

	return name, userID, password, nil
}
