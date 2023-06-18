package postgres

import (
	"context"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchAllLessons(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// container and database
	container, db, err := createPostgresContainer(ctx, "testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer container.Terminate(ctx)

	// migration
	mig, err := newMigrator(db)
	if err != nil {
		t.Fatal(err)
	}

	err = mig.Up()
	if err != nil {
		t.Fatal(err)
	}

	r := NewReader(db)

	lesson := &Lesson{}
	createLesson(db, lesson)

	res, err := r.FetchAllLessons()

	assert.NoError(t, err)
	assert.IsType(t, models.Lessons{}, res)
	assert.Len(t, res, 1)
	assert.Equal(t, lesson.Id, res[0].Id)
	assert.Equal(t, lesson.Name, res[0].Name)
	assert.Equal(t, lesson.Description, res[0].Description)
}

func TestHydrateLesson(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// container and database
	container, db, err := createPostgresContainer(ctx, "testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer container.Terminate(ctx)

	// migration
	mig, err := newMigrator(db)
	if err != nil {
		t.Fatal(err)
	}

	err = mig.Up()
	if err != nil {
		t.Fatal(err)
	}

	r := NewReader(db)

	l := &Lesson{
		Name:          "foo",
		Description:   "bar",
		ExerciseCount: 10,
	}
	createLesson(db, l)

	lesson := &models.Lesson{
		Id: l.Id,
	}

	err = r.HydrateLesson(lesson)

	assert.NoError(t, err)
	assert.Equal(t, lesson.Name, l.Name)
	assert.Equal(t, lesson.Description, l.Description)
	assert.Equal(t, lesson.ExerciseCount, l.ExerciseCount)
}

func TestFetchExercisesOfLesson(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// container and database
	container, db, err := createPostgresContainer(ctx, "testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer container.Terminate(ctx)

	// migration
	mig, err := newMigrator(db)
	if err != nil {
		t.Fatal(err)
	}

	err = mig.Up()
	if err != nil {
		t.Fatal(err)
	}

	r := NewReader(db)

	exercise := &Exercise{}
	createExercise(db, exercise)

	// to check of exercise without answers will also be fetched
	exercise2 := &Exercise{LessonId: exercise.LessonId}
	createExercise(db, exercise2)

	// belongs to exercise, to be included
	answer1 := &Answer{ExerciseId: exercise.Id}
	createAnswer(db, answer1)

	// does not belong to exercise, to be excluded
	answer2 := &Answer{}
	createAnswer(db, answer2)

	res, err := r.FetchExercisesOfLesson(models.Lesson{Id: exercise.LessonId})

	assert.NoError(t, err)
	assert.IsType(t, models.Exercises{}, res)
	assert.Len(t, res, 2)
	assert.Equal(t, exercise.Id, res[1].Id)
	assert.Equal(t, exercise.Question, res[1].Question)
	assert.Equal(t, exercise.Answer, res[1].Answer)
	assert.Len(t, res[1].Answers, 1)
	assert.Empty(t, res[0].Answers)
	assert.Equal(t, answer1.Id, res[1].Answers[0].Id)
	assert.Equal(t, answer1.Type, res[1].Answers[0].Type)
	assert.Equal(t, answer1.CreatedAt, res[1].Answers[0].CreatedAt)
}

func TestFetchAnswersOfExercise(t *testing.T) {
	t.Parallel()

	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	ctx := context.Background()

	// container and database
	container, db, err := createPostgresContainer(ctx, "testdb")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer container.Terminate(ctx)

	// migration
	mig, err := newMigrator(db)
	if err != nil {
		t.Fatal(err)
	}

	err = mig.Up()
	if err != nil {
		t.Fatal(err)
	}

	r := NewReader(db)

	exercise := &Exercise{}
	createExercise(db, exercise)

	// store good answer related to the exercise
	answer := &Answer{
		ExerciseId: exercise.Id,
	}
	createAnswer(db, answer)

	// store bad answer related to another exercise
	createAnswer(db, &Answer{})

	answers, err := r.FetchAnswersOfExercise(models.Exercise{Id: exercise.Id})

	assert.NoError(t, err)
	assert.Len(t, answers, 1)
	assert.Equal(t, answer.Id, answers[0].Id)
}
