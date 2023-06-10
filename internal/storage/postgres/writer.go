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

func (w *Writer) StoreLesson(lesson *models.Lesson) error {
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

func (w *Writer) StoreExercise(exercise *models.Exercise) error {
	var query string

	if exercise.Id > 0 {
		query = `UPDATE exercise set question = $1, answer = $2 where id = $3;`

		_, err := w.db.Exec(query, exercise.Question, exercise.Answer, exercise.Id)
		if err != nil {
			return fmt.Errorf("failed to execute 'UPDATE exercise' query: %w", err)
		}
	} else {
		// Get a Tx for making transaction requests.
		tx, err := w.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to begin DB transaction: %w", err)
		}

		// Defer a rollback in case anything fails.
		defer tx.Rollback()

		// Insert exercise.

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

		// Update lesson.exercise_count.

		query = `
			UPDATE lesson
			SET exercise_count=sq.exercise_count
			FROM (SELECT count(*) as exercise_count FROM exercise WHERE lesson_id = $1) AS sq
			WHERE lesson.id=$1;
			`

		_, err = tx.Exec(query, exercise.Lesson.Id)
		if err != nil {
			return fmt.Errorf("failed to execute 'UPDATE lesson' query: %w", err)
		}

		// Commit the transaction.
		if err = tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit DB transaction: %w", err)
		}
	}

	return nil
}

func (w *Writer) DeleteExercise(exercise models.Exercise) error {
	var query string

	// Get a Tx for making transaction requests.
	tx, err := w.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin DB transaction: %w", err)
	}

	// Defer a rollback in case anything fails.
	defer tx.Rollback()

	// Fetch the lesson of the exercise, so it can be used to update lesson.exercise_count.

	query = `SELECT lesson_id FROM exercise WHERE id = $1;`

	exercise.Lesson = &models.Lesson{}

	if err := tx.QueryRow(query, exercise.Id).Scan(&exercise.Lesson.Id); err != nil {
		return fmt.Errorf("failed to execute 'SELECT lesson_id FROM exercise' query: %w", err)
	}

	// Delete exercise.

	query = `DELETE FROM exercise WHERE id = $1;`

	_, err = tx.Exec(query, exercise.Id)
	if err != nil {
		return fmt.Errorf("failed to execute 'DELETE FROM exercise' query: %w", err)
	}

	// Update lesson.exercise_count.

	query = `
			UPDATE lesson
			SET exercise_count=sq.exercise_count
			FROM (SELECT count(*) as exercise_count FROM exercise WHERE lesson_id = $1) AS sq
			WHERE lesson.id=$1;
			`

	_, err = tx.Exec(query, exercise.Lesson.Id)
	if err != nil {
		return fmt.Errorf("failed to execute 'UPDATE lesson' query: %w", err)
	}

	// Commit the transaction.
	if err = tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit DB transaction: %w", err)
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
