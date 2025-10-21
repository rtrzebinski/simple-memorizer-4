package postgres

import (
	"database/sql"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

type (
	lesson struct {
		Id          int
		Name        string
		Description string
		UserID      int
	}

	exercise struct {
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

	result struct {
		Id         int
		ExerciseId int
		Type       string
		CreatedAt  time.Time
	}

	user struct {
		Id       int
		Name     string
		Email    string
		Password string
	}
)

func createLesson(db *sql.DB, lesson *lesson) {
	query := `INSERT INTO lesson (name, description, user_id) VALUES ($1, $2, $3) RETURNING id;`

	if lesson.Name == "" {
		lesson.Name = randomString()
	}

	if lesson.Description == "" {
		lesson.Description = randomString()
	}

	if lesson.UserID == 0 {
		user := user{}
		user.Name = randomString()
		user.Email = randomString()
		user.Password = randomString()
		createUser(db, &user)
		lesson.UserID = user.Id
	}

	err := db.QueryRow(query, &lesson.Name, &lesson.Description, &lesson.UserID).Scan(&lesson.Id)
	if err != nil {
		panic(err)
	}
}

func findLessonById(db *sql.DB, lessonId int) *lesson {
	const query = `
		SELECT l.id, l.name, l.description
		FROM lesson l
		WHERE l.id = $1;`

	rows, err := db.Query(query, lessonId)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var lesson lesson

		err = rows.Scan(&lesson.Id, &lesson.Name, &lesson.Description)
		if err != nil {
			panic(err)
		}

		return &lesson
	}

	return nil
}

func fetchLatestLesson(db *sql.DB) lesson {
	var lesson lesson

	const query = `
		SELECT l.id, l.name, l.description, l.user_id
		FROM lesson l
		ORDER BY id DESC
		LIMIT 1;`

	if err := db.QueryRow(query).Scan(&lesson.Id, &lesson.Name, &lesson.Description, &lesson.UserID); err != nil {
		panic(err)
	}

	return lesson
}

func createExercise(db *sql.DB, exercise *exercise) {
	query := `
INSERT INTO exercise (lesson_id, question, answer, bad_answers, bad_answers_today, latest_bad_answer, 
latest_bad_answer_was_today, good_answers, good_answers_today, latest_good_answer, latest_good_answer_was_today)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;
`

	if exercise.LessonId == 0 {
		lesson := lesson{}
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

func findExerciseById(db *sql.DB, exerciseId int) *exercise {
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
		var exercise exercise

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

func fetchLatestExercise(db *sql.DB) exercise {
	var exercise exercise

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

func fetchLatestResult(db *sql.DB) result {
	var answer result

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

func createResult(db *sql.DB, answer *result) {
	query := `INSERT INTO result (exercise_id, type) VALUES ($1, $2) RETURNING id, created_at;`

	if answer.ExerciseId == 0 {
		exercise := &exercise{}
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

func createUser(db *sql.DB, user *user) {
	query := `INSERT INTO "user" (name, email, password) VALUES ($1, $2, $3) RETURNING id;`

	err := db.QueryRow(query, &user.Name, &user.Email, &user.Password).
		Scan(&user.Id)
	if err != nil {
		panic(err)
	}
}

func fetchUserByEmail(db *sql.DB, email string) *user {
	const query = `SELECT id, name, email, password FROM "user" WHERE email = $1;`

	rows, err := db.Query(query, email)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var user user

		err = rows.Scan(&user.Id, &user.Name, &user.Email, &user.Password)
		if err != nil {
			panic(err)
		}

		return &user
	}

	return nil
}

func randomString() string {
	return uuid.NewString()
}
