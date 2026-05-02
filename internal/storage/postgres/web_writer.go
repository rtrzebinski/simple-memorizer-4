package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type WebWriter struct {
	db *sql.DB
}

func NewWebWriter(db *sql.DB) *WebWriter {
	return &WebWriter{db: db}
}

func (w *WebWriter) UpsertLesson(ctx context.Context, userID string, lesson *backend.Lesson) error {
	slog.Debug("WebWriter UpsertLesson", "userID", userID)

	var query string

	if lesson.Id > 0 {
		query = `UPDATE lesson SET name = $1, description = $2 WHERE id = $3 AND user_id = $4;`

		_, err := w.db.ExecContext(ctx, query, lesson.Name, lesson.Description, lesson.Id, userID)
		if err != nil {
			return fmt.Errorf("failed to update lesson: %w", err)
		}
	} else {
		query = `INSERT INTO lesson (name, description, user_id) VALUES ($1, $2, $3) RETURNING id;`

		rows, err := w.db.QueryContext(ctx, query, lesson.Name, lesson.Description, userID)
		if err != nil {
			return fmt.Errorf("failed to insert lesson: %w", err)
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

func (w *WebWriter) DeleteLesson(ctx context.Context, userID string, lesson backend.Lesson) error {
	slog.Debug("WebWriter DeleteLesson", "userID", userID)

	query := `DELETE FROM lesson WHERE id = $1 AND user_id = $2;`

	_, err := w.db.ExecContext(ctx, query, lesson.Id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete lesson: %w", err)
	}

	return nil
}

func (w *WebWriter) UpsertExercise(ctx context.Context, userID string, exercise *backend.Exercise) error {
	slog.Debug("WebWriter UpsertExercise", "userID", userID)

	var query string

	if exercise.Id > 0 {
		query = `
			UPDATE exercise e 
			SET question = $1, answer = $2 
			FROM lesson l 
			WHERE e.lesson_id = l.id 
	  		AND e.id = $3 
	  		AND l.user_id = $4;
		`

		_, err := w.db.ExecContext(ctx, query, exercise.Question, exercise.Answer, exercise.Id, userID)
		if err != nil {
			return fmt.Errorf("failed to update exercise: %w", err)
		}
	} else {
		query = `
			INSERT INTO exercise (lesson_id, question, answer)
			SELECT $1, $2, $3
			FROM lesson l
			WHERE l.id = $1 AND l.user_id = $4
			RETURNING id;
		`

		rows, err := w.db.QueryContext(ctx, query, exercise.Lesson.Id, exercise.Question, exercise.Answer, userID)
		if err != nil {
			return fmt.Errorf("failed to insert exercise: %w", err)
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

func (w *WebWriter) StoreExercises(ctx context.Context, userID string, exercises backend.Exercises) error {
	slog.Debug("WebWriter StoreExercises", "userID", userID)

	const query = `
		INSERT INTO exercise (lesson_id, question, answer)
		SELECT $1, $2, $3
		FROM lesson l
		WHERE l.id = $1 AND l.user_id = $4
		ON CONFLICT (lesson_id, question) DO NOTHING;
	`

	for _, exercise := range exercises {
		_, err := w.db.ExecContext(ctx, query, exercise.Lesson.Id, exercise.Question, exercise.Answer, userID)
		if err != nil {
			return fmt.Errorf("failed to store exercise: %w", err)
		}
	}

	return nil
}

func (w *WebWriter) DeleteExercise(ctx context.Context, userID string, exercise backend.Exercise) error {
	slog.Debug("WebWriter DeleteExercise", "userID", userID)

	query := `
		DELETE FROM exercise e
		USING lesson l
		WHERE e.lesson_id = l.id
	  	AND e.id = $1
	  	AND l.user_id = $2;
	`

	_, err := w.db.ExecContext(ctx, query, exercise.Id, userID)
	if err != nil {
		return fmt.Errorf("failed to delete exercise: %w", err)
	}

	return nil
}
