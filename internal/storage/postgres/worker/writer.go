package worker

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/worker"
)

type Writer struct {
	db *sql.DB
}

func NewWriter(db *sql.DB) *Writer {
	return &Writer{db: db}
}

func (w *Writer) StoreResult(ctx context.Context, result worker.Result) error {
	const query = `INSERT INTO result (type, exercise_id) VALUES ($1, $2) RETURNING id;`

	rows, err := w.db.QueryContext(ctx, query, result.Type, result.ExerciseId)
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

func (w *Writer) UpdateExerciseProjection(ctx context.Context, exerciseID int, projection worker.ResultsProjection) error {
	const query = `UPDATE exercise SET
                    bad_answers = $1,
                    bad_answers_today = $2,
                    latest_bad_answer = $3,
                    latest_bad_answer_was_today = $4,
                    good_answers = $5,
                    good_answers_today = $6,
                    latest_good_answer = $7,
                    latest_good_answer_was_today = $8
                    WHERE id = $9;`

	_, err := w.db.ExecContext(ctx, query, projection.BadAnswers, projection.BadAnswersToday, projection.LatestBadAnswer,
		projection.LatestBadAnswerWasToday, projection.GoodAnswers, projection.GoodAnswersToday,
		projection.LatestGoodAnswer, projection.LatestGoodAnswerWasToday, exerciseID)
	if err != nil {
		return fmt.Errorf("execute 'UPDATE exercise' query: %w", err)
	}

	return nil
}
