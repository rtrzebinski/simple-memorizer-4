package postgres

import (
	"database/sql"
	"github.com/rtrzebinski/simple-memorizer-go/internal/models"
)

type Reader struct {
	db *sql.DB
}

func NewReader(db *sql.DB) *Reader {
	return &Reader{db: db}
}

func (r *Reader) RandomExercise() models.Exercise {
	var exercise models.Exercise

	const query = `SELECT id, question, answer FROM exercise ORDER BY random() LIMIT 1`

	if err := r.db.QueryRow(query).Scan(&exercise.Id, &exercise.Question, &exercise.Answer); err != nil {
		panic(err)
	}

	return exercise
}
