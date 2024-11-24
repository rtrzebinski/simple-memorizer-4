package postgres

import (
	"context"
	"database/sql"
	"github.com/rtrzebinski/simple-memorizer-4/internal/worker"
)

type Reader struct {
	db *sql.DB
}

func NewReader(db *sql.DB) *Reader {
	return &Reader{db: db}
}

func (r *Reader) FetchResults(ctx context.Context, exerciseID int) ([]worker.Result, error) {
	const query = `SELECT id, type, exercise_id, created_at FROM result WHERE exercise_id = $1;`

	rows, err := r.db.QueryContext(ctx, query, exerciseID)
	if err != nil {
		return nil, err
	}

	var results []worker.Result
	for rows.Next() {
		var result worker.Result
		err = rows.Scan(&result.Id, &result.Type, &result.ExerciseId, &result.CreatedAt)
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}
