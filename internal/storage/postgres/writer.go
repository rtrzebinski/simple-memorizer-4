package postgres

import (
	"database/sql"
	"log"
)

type Writer struct {
	db *sql.DB
}

func NewWriter(db *sql.DB) *Writer {
	return &Writer{db: db}
}

// todo add unique key to exercise_result.exercise_id

func (w *Writer) IncrementBadAnswers(exerciseId int) {
	// check for existing exercise result
	query := `SELECT id FROM exercise_result where exercise_id = $1`

	var exerciseResultId int

	err := w.db.QueryRow(query, exerciseId).Scan(&exerciseResultId)

	// exercise result does not exist - create it
	if err != nil && err == sql.ErrNoRows {
		query = `INSERT INTO exercise_result (exercise_id, bad_answers) VALUES ($1, 1);`

		_, err := w.db.Exec(query, exerciseId)
		if err != nil {
			panic(err)
		}

		log.Println("Created exercise_result for bad_answer")

		return
	}

	// exercise result exist - increment bad_answers
	query = `UPDATE exercise_result SET bad_answers = bad_answers + 1 WHERE exercise_id = $1;`

	_, err = w.db.Exec(query, exerciseId)
	if err != nil {
		panic(err)
	}

	log.Println("Incremented exercise_result bad_answers")
}

func (w *Writer) IncrementGoodAnswers(exerciseId int) {
	// check for existing exercise result
	query := `SELECT id FROM exercise_result where exercise_id = $1`

	var exerciseResultId int

	err := w.db.QueryRow(query, exerciseId).Scan(&exerciseResultId)

	// exercise result does not exist - create it
	if err != nil && err == sql.ErrNoRows {
		query = `INSERT INTO exercise_result (exercise_id, good_answers) VALUES ($1, 1);`

		_, err := w.db.Exec(query, exerciseId)
		if err != nil {
			panic(err)
		}

		log.Println("Created exercise_result for good_answer")

		return
	}

	// exercise result exist - increment good_answers
	query = `UPDATE exercise_result SET good_answers = good_answers + 1 WHERE exercise_id = $1;`

	_, err = w.db.Exec(query, exerciseId)
	if err != nil {
		panic(err)
	}

	log.Println("Incremented exercise_result good_answers")
}
