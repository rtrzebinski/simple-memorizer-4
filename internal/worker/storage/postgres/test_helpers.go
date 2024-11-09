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

func randomString() string {
	return uuid.NewString()
}
