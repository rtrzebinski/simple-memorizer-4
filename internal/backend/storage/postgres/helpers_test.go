package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/docker/go-connections/nat"
	"github.com/golang-migrate/migrate/v4"
	migrate_postgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"time"
)

type (
	Lesson struct {
		Id            int
		Name          string
		ExerciseCount int
	}

	Exercise struct {
		Id       int
		LessonId int
		Question string
		Answer   string
	}

	ExerciseResult struct {
		Id          int
		ExerciseId  int
		BadAnswers  int
		GoodAnswers int
	}
)

func findExerciseById(db *sql.DB, exerciseId int) *Exercise {
	const query = `
		SELECT e.id, e.question, e.answer
		FROM exercise e
		WHERE e.id = $1;`

	rows, err := db.Query(query, exerciseId)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var exercise Exercise

		err = rows.Scan(&exercise.Id, &exercise.Question, &exercise.Answer)
		if err != nil {
			panic(err)
		}

		return &exercise
	}

	return nil
}

func findLessonById(db *sql.DB, lessonId int) *Lesson {
	const query = `
		SELECT l.id, l.name, l.exercise_count
		FROM lesson l
		WHERE l.id = $1;`

	rows, err := db.Query(query, lessonId)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var lesson Lesson

		err = rows.Scan(&lesson.Id, &lesson.Name, &lesson.ExerciseCount)
		if err != nil {
			panic(err)
		}

		return &lesson
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

func fetchLatestLesson(db *sql.DB) Lesson {
	var lesson Lesson

	const query = `
		SELECT l.id, l.name
		FROM lesson l
		ORDER BY id DESC
		LIMIT 1;`

	if err := db.QueryRow(query).Scan(&lesson.Id, &lesson.Name); err != nil {
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

func createLesson(db *sql.DB, lesson *Lesson) {
	query := `INSERT INTO lesson (name, exercise_count) VALUES ($1, $2) RETURNING id;`

	if lesson.Name == "" {
		lesson.Name = randomString()
	}

	err := db.QueryRow(query, &lesson.Name, &lesson.ExerciseCount).Scan(&lesson.Id)
	if err != nil {
		panic(err)
	}
}

func randomString() string {
	return uuid.NewString()
}

func findExerciseResultByExerciseId(db *sql.DB, exerciseId int) ExerciseResult {
	var exerciseResult ExerciseResult

	const query = `
		SELECT er.id, er.exercise_id, er.bad_answers, er.good_answers
		FROM exercise_result er
		WHERE er.exercise_id = $1;`

	if err := db.QueryRow(query, exerciseId).Scan(&exerciseResult.Id, &exerciseResult.ExerciseId, &exerciseResult.BadAnswers, &exerciseResult.GoodAnswers); err != nil {
		panic(err)
	}

	return exerciseResult
}

func createPostgresContainer(ctx context.Context, dbname string) (testcontainers.Container, *sql.DB, error) {
	env := map[string]string{
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_USER":     "postgres",
		"POSTGRES_DB":       dbname,
	}

	port := "5432/tcp"
	dbURL := func(host string, port nat.Port) string {
		return fmt.Sprintf("postgres://postgres:password@%s:%s/%s?sslmode=disable", host, port.Port(), dbname)
	}

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:15-alpine",
			ExposedPorts: []string{port},
			Cmd:          []string{"postgres", "-c", "fsync=off"},
			Env:          env,
			SkipReaper:   true,
			WaitingFor: wait.ForAll(
				wait.ForSQL(nat.Port(port), "postgres", dbURL).WithStartupTimeout(time.Second*5),
				wait.ForLog("database system is ready to accept connections"),
			),
		},
		Started: true,
	}

	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, fmt.Errorf("failed to start container: %s", err)
	}

	mappedPort, err := container.MappedPort(ctx, nat.Port(port))
	if err != nil {
		return container, nil, fmt.Errorf("failed to get container external port: %s", err)
	}

	log.Println("postgres container ready and running at port: ", mappedPort)

	url := fmt.Sprintf("postgres://postgres:password@localhost:%s/%s?sslmode=disable", mappedPort.Port(), dbname)

	db, err := sql.Open("postgres", url)
	if err != nil {
		return container, db, fmt.Errorf("failed to establish database connection: %s", err)
	}

	return container, db, nil
}

func newMigrator(db *sql.DB) (*migrate.Migrate, error) {
	driver, err := migrate_postgres.WithInstance(db, &migrate_postgres.Config{})
	if err != nil {
		log.Fatalf("failed to create migrator driver: %s", err)
	}

	return migrate.NewWithDatabaseInstance("file://../../../../migrations", "postgres", driver)
}
