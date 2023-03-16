package postgres

import (
	"database/sql"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

type Reader struct {
	db *sql.DB
}

func NewReader(db *sql.DB) *Reader {
	return &Reader{db: db}
}

func (r *Reader) RandomExercise() models.Exercise {
	var exercise models.Exercise

	const query = `
		SELECT e.id, e.question, e.answer, COALESCE(er.bad_answers, 0), COALESCE(er.good_answers, 0) 
		FROM exercise e
		LEFT JOIN exercise_result er on e.id = er.exercise_id
		ORDER BY random()
		LIMIT 1`

	if err := r.db.QueryRow(query).Scan(&exercise.Id, &exercise.Question, &exercise.Answer, &exercise.BadAnswers, &exercise.GoodAnswers); err != nil {
		panic(err)
	}

	return exercise
}
