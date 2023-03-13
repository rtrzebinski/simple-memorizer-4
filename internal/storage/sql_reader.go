package storage

import (
	"database/sql"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-go/internal/models"
)

type SqlReader struct {
	db *sql.DB
}

func NewSqlReader(db *sql.DB) *SqlReader {
	return &SqlReader{db: db}
}

func (r *SqlReader) RandomExercise() models.Exercise {
	var exercise models.Exercise

	const query = `SELECT question, answer FROM exercise ORDER BY random() LIMIT 1`

	if err := r.db.QueryRow(query).Scan(&exercise.Question, &exercise.Answer); err != nil {
		fmt.Println(err)
	}

	return exercise
}
