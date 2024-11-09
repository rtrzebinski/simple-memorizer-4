package postgres

import (
	"database/sql"
	"fmt"

	"github.com/rtrzebinski/simple-memorizer-4/internal/worker"
)

type Writer struct {
	db *sql.DB
}

func NewWriter(db *sql.DB) *Writer {
	return &Writer{db: db}
}

func (w *Writer) StoreResult(result *worker.Result) error {
	var query string

	query = `INSERT INTO result (type, exercise_id) VALUES ($1, $2) RETURNING id;`

	rows, err := w.db.Query(query, result.Type, result.ExerciseId)
	if err != nil {
		return fmt.Errorf("failed to execute 'INSERT INTO result' query: %w", err)
	}

	for rows.Next() {
		err = rows.Scan(&result.Id)
		if err != nil {
			return fmt.Errorf("failed to scan result insert id: %w", err)
		}
	}

	return nil
}
