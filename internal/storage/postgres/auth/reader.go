package auth

import (
	"context"
	"database/sql"
)

type Reader struct {
	db *sql.DB
}

func NewReader(db *sql.DB) *Reader {
	return &Reader{db: db}
}

func (r *Reader) FetchUser(ctx context.Context, email string) (name, userID, password string, err error) {
	const query = `SELECT name, id, password FROM "user" WHERE email = $1;`

	row := r.db.QueryRowContext(ctx, query, email)

	err = row.Scan(&name, &userID, &password)
	if err != nil {
		return "", "", "", err
	}

	return name, userID, password, nil
}
