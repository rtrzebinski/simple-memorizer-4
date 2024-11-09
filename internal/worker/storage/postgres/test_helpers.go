package postgres

import (
	"time"

	"database/sql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
)

// Tables

type (
	Lesson struct {
		Id          int
		Name        string
		Description string
	}

	Exercise struct {
		Id                       int
		LessonId                 int
		Question                 string
		Answer                   string
		BadAnswers               int
		BadAnswersToday          int
		LatestBadAnswer          time.Time
		LatestBadAnswerWasToday  bool
		GoodAnswers              int
		GoodAnswersToday         int
		LatestGoodAnswer         time.Time
		LatestGoodAnswerWasToday bool
	}

	Result struct {
		Id         int
		ExerciseId int
		Type       string
		CreatedAt  time.Time
	}
)

func createLesson(db *sql.DB, lesson *Lesson) {
	query := `INSERT INTO lesson (name, description) VALUES ($1, $2) RETURNING id;`

	if lesson.Name == "" {
		lesson.Name = randomString()
	}

	if lesson.Description == "" {
		lesson.Description = randomString()
	}

	err := db.QueryRow(query, &lesson.Name, &lesson.Description).Scan(&lesson.Id)
	if err != nil {
		panic(err)
	}
}

func createExercise(db *sql.DB, exercise *Exercise) {
	query := `INSERT INTO exercise (lesson_id, question, answer) VALUES ($1, $2, $3) RETURNING id;`

	if exercise.LessonId == 0 {
		lesson := Lesson{}
		createLesson(db, &lesson)
		exercise.LessonId = lesson.Id
	}

	if exercise.Question == "" {
		exercise.Question = randomString()
	}

	if exercise.Answer == "" {
		exercise.Answer = randomString()
	}

	err := db.QueryRow(query, &exercise.LessonId, &exercise.Question, &exercise.Answer).Scan(&exercise.Id)
	if err != nil {
		panic(err)
	}
}

func findExerciseById(db *sql.DB, exerciseId int) *Exercise {
	const query = `
		SELECT e.id, e.question, e.answer, e.lesson_id, e.bad_answers, e.bad_answers_today, e.latest_bad_answer,
		e.latest_bad_answer_was_today, e.good_answers, e.good_answers_today, e.latest_good_answer,
		e.latest_good_answer_was_today
		FROM exercise e
		WHERE e.id = $1;`

	rows, err := db.Query(query, exerciseId)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var exercise Exercise

		err = rows.Scan(&exercise.Id, &exercise.Question, &exercise.Answer, &exercise.LessonId, &exercise.BadAnswers,
			&exercise.BadAnswersToday, &exercise.LatestBadAnswer, &exercise.LatestBadAnswerWasToday,
			&exercise.GoodAnswers, &exercise.GoodAnswersToday, &exercise.LatestGoodAnswer,
			&exercise.LatestGoodAnswerWasToday)
		if err != nil {
			panic(err)
		}

		return &exercise
	}

	return nil
}

func fetchLatestResult(db *sql.DB) Result {
	var answer Result

	const query = `
		SELECT r.id, r.type, r.exercise_id, r.created_at
		FROM result r
		ORDER BY id DESC
		LIMIT 1;`

	if err := db.QueryRow(query).Scan(&answer.Id, &answer.Type, &answer.ExerciseId, &answer.CreatedAt); err != nil {
		panic(err)
	}

	return answer
}

func createResult(db *sql.DB, answer *Result) {
	query := `INSERT INTO result (exercise_id, type) VALUES ($1, $2) RETURNING id, created_at;`

	if answer.ExerciseId == 0 {
		exercise := &Exercise{}
		createExercise(db, exercise)
		answer.ExerciseId = exercise.Id
	}

	if answer.Type == "" {
		answer.Type = "good"
	}

	err := db.QueryRow(query, &answer.ExerciseId, &answer.Type).Scan(&answer.Id, &answer.CreatedAt)
	if err != nil {
		panic(err)
	}
}

func randomString() string {
	return uuid.NewString()
}
