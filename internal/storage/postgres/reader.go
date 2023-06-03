package postgres

import (
	"database/sql"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

type Reader struct {
	db *sql.DB
}

func NewReader(db *sql.DB) *Reader {
	return &Reader{db: db}
}

func (r *Reader) FetchAllLessons() (models.Lessons, error) {
	var lessons models.Lessons

	const query = `
		SELECT l.id, l.name, l.exercise_count
		FROM lesson l
		ORDER BY l.id DESC
		`

	rows, err := r.db.Query(query)
	if err != nil {
		return lessons, err
	}

	for rows.Next() {
		var lesson models.Lesson

		err = rows.Scan(&lesson.Id, &lesson.Name, &lesson.ExerciseCount)
		if err != nil {
			return lessons, err
		}

		lessons = append(lessons, lesson)
	}

	return lessons, nil
}

func (r *Reader) FetchExercisesOfLesson(lesson models.Lesson) (models.Exercises, error) {
	var exercises models.Exercises

	const query = `
		SELECT e.id, e.question, e.answer, COALESCE(er.bad_answers, 0), COALESCE(er.good_answers, 0) 
		FROM exercise e
		LEFT JOIN exercise_result er on e.id = er.exercise_id
		WHERE lesson_id = $1
		ORDER BY e.id DESC
		`

	rows, err := r.db.Query(query, lesson.Id)
	if err != nil {
		return exercises, err
	}

	for rows.Next() {
		var exercise models.Exercise

		err = rows.Scan(&exercise.Id, &exercise.Question, &exercise.Answer, &exercise.BadAnswers, &exercise.GoodAnswers)
		if err != nil {
			return exercises, err
		}

		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func (r *Reader) FetchRandomExerciseOfLesson(lesson models.Lesson) (models.Exercise, error) {
	var exercise models.Exercise

	const query = `
		SELECT e.id, e.question, e.answer, COALESCE(er.bad_answers, 0), COALESCE(er.good_answers, 0) 
		FROM exercise e
		JOIN lesson l on l.id = e.lesson_id
		LEFT JOIN exercise_result er on e.id = er.exercise_id
		WHERE l.id = $1
		ORDER BY random()
		LIMIT 1`

	if err := r.db.QueryRow(query, lesson.Id).Scan(&exercise.Id, &exercise.Question, &exercise.Answer, &exercise.BadAnswers, &exercise.GoodAnswers); err != nil {
		return exercise, fmt.Errorf("failed to scan query results: %w", err)
	}

	return exercise, nil
}
