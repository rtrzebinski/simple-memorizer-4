package web

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type Writer struct {
	db *sql.DB
}

func NewWriter(db *sql.DB) *Writer {
	return &Writer{db: db}
}

func (w *Writer) UpsertLesson(ctx context.Context, lesson *backend.Lesson, userID string) error {
	slog.Debug("Writer UpsertLesson", "userID", userID)

	var query string

	if lesson.Id > 0 {
		query = `UPDATE lesson set name = $1, description = $2 where id = $3;`

		_, err := w.db.ExecContext(ctx, query, lesson.Name, lesson.Description, lesson.Id)
		if err != nil {
			return fmt.Errorf("failed to execute 'UPDATE lesson' query: %w", err)
		}
	} else {
		query = `INSERT INTO lesson (name, description, user_id) VALUES ($1, $2, $3) RETURNING id;`

		rows, err := w.db.QueryContext(ctx, query, lesson.Name, lesson.Description, userID)
		if err != nil {
			return fmt.Errorf("failed to execute 'INSERT INTO lesson' query: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&lesson.Id)
			if err != nil {
				return fmt.Errorf("failed to scan lesson insert id: %w", err)
			}
		}
	}

	return nil
}

func (w *Writer) DeleteLesson(ctx context.Context, lesson backend.Lesson, userID string) error {
	slog.Debug("Writer DeleteLesson", "userID", userID)

	query := `DELETE FROM lesson WHERE id = $1;`

	_, err := w.db.ExecContext(ctx, query, lesson.Id)
	if err != nil {
		return fmt.Errorf("failed to execute 'DELETE FROM lesson' query: %w", err)
	}

	return nil
}

func (w *Writer) UpsertExercise(ctx context.Context, exercise *backend.Exercise, userID string) error {
	slog.Debug("Writer UpsertExercise", "userID", userID)

	var query string

	if exercise.Id > 0 {
		query = `UPDATE exercise set question = $1, answer = $2 where id = $3;`

		_, err := w.db.ExecContext(ctx, query, exercise.Question, exercise.Answer, exercise.Id)
		if err != nil {
			return fmt.Errorf("failed to execute 'UPDATE exercise' query: %w", err)
		}
	} else {
		query = `INSERT INTO exercise (lesson_id, question, answer) VALUES ($1, $2, $3) RETURNING id;`

		rows, err := w.db.QueryContext(ctx, query, exercise.Lesson.Id, exercise.Question, exercise.Answer)
		if err != nil {
			return fmt.Errorf("failed to execute 'INSERT INTO exercise' query: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&exercise.Id)
			if err != nil {
				return fmt.Errorf("failed to scan exercise insert id: %w", err)
			}
		}
	}

	return nil
}

func (w *Writer) StoreExercises(ctx context.Context, exercises backend.Exercises, userID string) error {
	slog.Debug("Writer StoreExercises", "userID", userID)

	const query = `
		INSERT INTO exercise (lesson_id, question, answer)
		VALUES ($1, $2, $3)
		ON CONFLICT (lesson_id, question) DO NOTHING;`

	for _, exercise := range exercises {
		_, err := w.db.ExecContext(ctx, query, exercise.Lesson.Id, exercise.Question, exercise.Answer)
		if err != nil {
			return fmt.Errorf("failed to execute 'INSERT INTO exercise' query: %w", err)
		}
	}

	return nil
}

func (w *Writer) DeleteExercise(ctx context.Context, exercise backend.Exercise, userID string) error {
	slog.Debug("Writer DeleteExercise", "userID", userID)

	query := `DELETE FROM exercise WHERE id = $1;`

	_, err := w.db.ExecContext(ctx, query, exercise.Id)
	if err != nil {
		return fmt.Errorf("failed to execute 'DELETE FROM exercise' query: %w", err)
	}

	return nil
}
