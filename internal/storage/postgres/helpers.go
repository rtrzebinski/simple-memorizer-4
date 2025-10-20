package postgres

import (
	"database/sql"
	"time"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/guregu/null/v5"
)

type (
	Lesson struct {
		Id          int
		Name        string
		Description string
		UserID      int
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

	Result struct {
		Id         int
		ExerciseId int
		Type       string
		CreatedAt  time.Time
	}

	User struct {
		Id       int
		Name     string
		Email    string
		Password string
	}
)

func CreateLesson(db *sql.DB, lesson *Lesson) {
	query := `INSERT INTO lesson (name, description, user_id) VALUES ($1, $2, $3) RETURNING id;`

	if lesson.Name == "" {
		lesson.Name = randomString()
	}

	if lesson.Description == "" {
		lesson.Description = randomString()
	}

	if lesson.UserID == 0 {
		user := User{}
		user.Name = randomString()
		user.Email = randomString()
		user.Password = randomString()
		CreateUser(db, &user)
		lesson.UserID = user.Id
	}

	err := db.QueryRow(query, &lesson.Name, &lesson.Description, &lesson.UserID).Scan(&lesson.Id)
	if err != nil {
		panic(err)
	}
}

func FindLessonById(db *sql.DB, lessonId int) *Lesson {
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

func FetchLatestLesson(db *sql.DB) Lesson {
	var lesson Lesson

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

func CreateExercise(db *sql.DB, exercise *Exercise) {
	query := `
INSERT INTO exercise (lesson_id, question, answer, bad_answers, bad_answers_today, latest_bad_answer, 
latest_bad_answer_was_today, good_answers, good_answers_today, latest_good_answer, latest_good_answer_was_today)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id;
`

	if exercise.LessonId == 0 {
		lesson := Lesson{}
		CreateLesson(db, &lesson)
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

func FindExerciseById(db *sql.DB, exerciseId int) *Exercise {
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

func FetchLatestExercise(db *sql.DB) Exercise {
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

func FetchLatestResult(db *sql.DB) Result {
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

func CreateResult(db *sql.DB, answer *Result) {
	query := `INSERT INTO result (exercise_id, type) VALUES ($1, $2) RETURNING id, created_at;`

	if answer.ExerciseId == 0 {
		exercise := &Exercise{}
		CreateExercise(db, exercise)
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

func CreateUser(db *sql.DB, user *User) {
	query := `INSERT INTO "user" (name, email, password) VALUES ($1, $2, $3) RETURNING id;`

	err := db.QueryRow(query, &user.Name, &user.Email, &user.Password).
		Scan(&user.Id)
	if err != nil {
		panic(err)
	}
}

func FetchUserByEmail(db *sql.DB, email string) *User {
	const query = `SELECT id, name, email, password FROM "user" WHERE email = $1;`

	rows, err := db.Query(query, email)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var user User

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
