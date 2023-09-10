package postgres

import (
	"database/sql"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

type Writer struct {
	db *sql.DB
}

func NewWriter(db *sql.DB) *Writer {
	return &Writer{db: db}
}

func (w *Writer) UpsertLesson(lesson *models.Lesson) error {
	var query string

	if lesson.Id > 0 {
		query = `UPDATE lesson set name = $1, description = $2 where id = $3;`

		_, err := w.db.Exec(query, lesson.Name, lesson.Description, lesson.Id)
		if err != nil {
			return fmt.Errorf("failed to execute 'UPDATE lesson' query: %w", err)
		}
	} else {
		query = `INSERT INTO lesson (name, description) VALUES ($1, $2) RETURNING id;`

		rows, err := w.db.Query(query, lesson.Name, lesson.Description)
		if err != nil {
			return fmt.Errorf("failed to execute 'INSERT INTO lesson' query: %w", err)
		}

		for rows.Next() {
			err = rows.Scan(&lesson.Id)
			if err != nil {
				return fmt.Errorf("failed to scan lesson insert id: %w", err)
			}
		}
	}

	return nil
}

func (w *Writer) DeleteLesson(lesson models.Lesson) error {
	query := `DELETE FROM lesson WHERE id = $1;`

	_, err := w.db.Exec(query, lesson.Id)
	if err != nil {
		return fmt.Errorf("failed to execute 'DELETE FROM lesson' query: %w", err)
	}

	return nil
}

func (w *Writer) UpsertExercise(exercise *models.Exercise) error {
	var query string

	if exercise.Id > 0 {
		query = `UPDATE exercise set question = $1, answer = $2 where id = $3;`

		_, err := w.db.Exec(query, exercise.Question, exercise.Answer, exercise.Id)
		if err != nil {
			return fmt.Errorf("failed to execute 'UPDATE exercise' query: %w", err)
		}
	} else {
		query = `INSERT INTO exercise (lesson_id, question, answer) VALUES ($1, $2, $3) RETURNING id;`

		rows, err := w.db.Query(query, exercise.Lesson.Id, exercise.Question, exercise.Answer)
		if err != nil {
			return fmt.Errorf("failed to execute 'INSERT INTO exercise' query: %w", err)
		}

		for rows.Next() {
			err = rows.Scan(&exercise.Id)
			if err != nil {
				return fmt.Errorf("failed to scan exercise insert id: %w", err)
			}
		}
	}

	return nil
}

func (w *Writer) StoreExercises(exercises models.Exercises) error {
	const query = `
		INSERT INTO exercise (lesson_id, question, answer)
		VALUES ($1, $2, $3)
		ON CONFLICT (lesson_id, question) DO NOTHING;`

	for _, exercise := range exercises {
		_, err := w.db.Query(query, exercise.Lesson.Id, exercise.Question, exercise.Answer)
		if err != nil {
			return fmt.Errorf("failed to execute 'INSERT INTO exercise' query: %w", err)
		}
	}

	return nil
}

func (w *Writer) DeleteExercise(exercise models.Exercise) error {

	query := `DELETE FROM exercise WHERE id = $1;`

	_, err := w.db.Exec(query, exercise.Id)
	if err != nil {
		return fmt.Errorf("failed to execute 'DELETE FROM exercise' query: %w", err)
	}

	return nil
}

func (w *Writer) StoreResult(result *models.Result) error {
	var query string

	query = `INSERT INTO result (type, exercise_id) VALUES ($1, $2) RETURNING id;`

	rows, err := w.db.Query(query, result.Type, result.Exercise.Id)
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
