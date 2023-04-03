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

func (w *Writer) DeleteExercise(exercise models.Exercise) error {
	query := `DELETE FROM exercise WHERE id = $1;`

	_, err := w.db.Exec(query, exercise.Id)
	if err != nil {
		return fmt.Errorf("failed to execute 'DELETE FROM exercise' query: %w", err)
	}

	return nil
}

func (w *Writer) StoreExercise(exercise models.Exercise) error {
	query := `INSERT INTO exercise (question, answer) VALUES ($1, $2);`

	_, err := w.db.Exec(query, exercise.Question, exercise.Answer)
	if err != nil {
		return fmt.Errorf("failed to execute 'INSERT INTO exercise' query: %w", err)
	}

	return nil
}

func (w *Writer) IncrementBadAnswers(exercise models.Exercise) error {
	// check for existing exercise result
	query := `SELECT id FROM exercise_result where exercise_id = $1`

	var exerciseResultId int

	err := w.db.QueryRow(query, exercise.Id).Scan(&exerciseResultId)

	// exercise result does not exist - create it
	if err != nil && err == sql.ErrNoRows {
		query = `INSERT INTO exercise_result (exercise_id, bad_answers) VALUES ($1, 1);`

		_, err := w.db.Exec(query, exercise.Id)
		if err != nil {
			return fmt.Errorf("failed to execute 'INSERT INTO exercise_result' query: %w", err)
		}

		return nil
	}

	// exercise result exist - increment bad_answers
	query = `UPDATE exercise_result SET bad_answers = bad_answers + 1 WHERE exercise_id = $1;`

	_, err = w.db.Exec(query, exercise.Id)
	if err != nil {
		return fmt.Errorf("failed to execute 'UPDATE exercise_result' query: %w", err)
	}

	return nil
}

func (w *Writer) IncrementGoodAnswers(exercise models.Exercise) error {
	// check for existing exercise result
	query := `SELECT id FROM exercise_result where exercise_id = $1`

	var exerciseResultId int

	err := w.db.QueryRow(query, exercise.Id).Scan(&exerciseResultId)

	// exercise result does not exist - create it
	if err != nil && err == sql.ErrNoRows {
		query = `INSERT INTO exercise_result (exercise_id, good_answers) VALUES ($1, 1);`

		_, err := w.db.Exec(query, exercise.Id)
		if err != nil {
			return fmt.Errorf("failed to execute 'INSERT INTO exercise_result' query: %w", err)
		}

		return nil
	}

	// exercise result exist - increment good_answers
	query = `UPDATE exercise_result SET good_answers = good_answers + 1 WHERE exercise_id = $1;`

	_, err = w.db.Exec(query, exercise.Id)
	if err != nil {
		return fmt.Errorf("failed to execute 'UPDATE exercise_result' query: %w", err)
	}

	return nil
}
