package postgres

import (
	"context"
	"database/sql"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/worker"
)

type WorkerReader struct {
	db *sql.DB
}

func NewWorkerReader(db *sql.DB) *WorkerReader {
	return &WorkerReader{db: db}
}

func (r *WorkerReader) FetchResults(ctx context.Context, exerciseID int) ([]worker.Result, error) {
	const query = `SELECT id, type, exercise_id, created_at FROM result WHERE exercise_id = $1;`

	rows, err := r.db.QueryContext(ctx, query, exerciseID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
