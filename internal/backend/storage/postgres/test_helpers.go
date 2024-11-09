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
		Id       int
		LessonId int
		Question string
		Answer   string
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
		SELECT e.id, e.question, e.answer, e.lesson_id
		FROM exercise e
		WHERE e.id = $1;`

	rows, err := db.Query(query, exerciseId)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var exercise Exercise

		err = rows.Scan(&exercise.Id, &exercise.Question, &exercise.Answer, &exercise.LessonId)
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
		SELECT e.id, e.lesson_id, e.question, e.answer
		FROM exercise e
		ORDER BY id DESC
		LIMIT 1;`

	if err := db.QueryRow(query).Scan(&exercise.Id, &exercise.LessonId, &exercise.Question, &exercise.Answer); err != nil {
		panic(err)
	}

	return exercise
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
