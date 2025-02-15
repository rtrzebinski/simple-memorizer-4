package postgres

import (
	"database/sql"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
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
		LatestBadAnswer          null.Time
		LatestBadAnswerWasToday  bool
		GoodAnswers              int
		GoodAnswersToday         int
		LatestGoodAnswer         null.Time
		LatestGoodAnswerWasToday bool
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

func findLessonById(db *sql.DB, lessonId int) *Lesson {
	const query = `
		SELECT l.id, l.name, l.description
		FROM lesson l
		WHERE l.id = $1;`

	rows, err := db.Query(query, lessonId)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var lesson Lesson

		err = rows.Scan(&lesson.Id, &lesson.Name, &lesson.Description)
		if err != nil {
			panic(err)
		}

		return &lesson
	}

	return nil
}

func fetchLatestLesson(db *sql.DB) Lesson {
	var lesson Lesson

	const query = `
		SELECT l.id, l.name, l.description
		FROM lesson l
		ORDER BY id DESC
		LIMIT 1;`

	if err := db.QueryRow(query).Scan(&lesson.Id, &lesson.Name, &lesson.Description); err != nil {
		panic(err)
	}

	return lesson
}

func createExercise(db *sql.DB, exercise *Exercise) {
	query := `
INSERT INTO exercise (lesson_id, question, answer, bad_answers, bad_answers_today, latest_bad_answer, 
latest_bad_answer_was_today, good_answers, good_answers_today, latest_good_answer, latest_good_answer_was_today)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;
`

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

	err := db.QueryRow(query, &exercise.LessonId, &exercise.Question, &exercise.Answer, &exercise.BadAnswers,
		&exercise.BadAnswersToday, &exercise.LatestBadAnswer, &exercise.LatestBadAnswerWasToday, &exercise.GoodAnswers,
		&exercise.GoodAnswersToday, &exercise.LatestGoodAnswer, &exercise.LatestGoodAnswerWasToday).Scan(&exercise.Id)
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

func fetchLatestExercise(db *sql.DB) Exercise {
	var exercise Exercise

	const query = `
SELECT e.id, e.lesson_id, e.question, e.answer, e.bad_answers, e.bad_answers_today, e.latest_bad_answer,
e.latest_bad_answer_was_today, e.good_answers, e.good_answers_today, e.latest_good_answer, e.latest_good_answer_was_today
FROM exercise e
ORDER BY id DESC
LIMIT 1;`

	if err := db.QueryRow(query).Scan(&exercise.Id, &exercise.LessonId, &exercise.Question, &exercise.Answer,
		&exercise.BadAnswers, &exercise.BadAnswersToday, &exercise.LatestBadAnswer, &exercise.LatestBadAnswerWasToday,
		&exercise.GoodAnswers, &exercise.GoodAnswersToday, &exercise.LatestGoodAnswer,
		&exercise.LatestGoodAnswerWasToday); err != nil {
		panic(err)
	}

	return exercise
}

func randomString() string {
	return uuid.NewString()
}
