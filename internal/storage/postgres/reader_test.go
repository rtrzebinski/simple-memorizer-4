package postgres

import (
	"context"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchLessons(t *testing.T) {
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

	createExercise(db, &Exercise{LessonId: lesson.Id})

	res, err := r.FetchLessons()

	assert.NoError(t, err)
	assert.IsType(t, models.Lessons{}, res)
	assert.Len(t, res, 1)
	assert.Equal(t, lesson.Id, res[0].Id)
	assert.Equal(t, lesson.Name, res[0].Name)
	assert.Equal(t, lesson.Description, res[0].Description)
	assert.Equal(t, 1, res[0].ExerciseCount)
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

	l := &Lesson{}
	createLesson(db, l)

	lesson := &models.Lesson{
		Id: l.Id,
	}

	err = r.HydrateLesson(lesson)

	assert.NoError(t, err)
	assert.Equal(t, l.Name, lesson.Name)
	assert.Equal(t, l.Description, lesson.Description)
	assert.Equal(t, 0, lesson.ExerciseCount)

	createExercise(db, &Exercise{LessonId: l.Id})
	createExercise(db, &Exercise{LessonId: l.Id})

	err = r.HydrateLesson(lesson)

	assert.NoError(t, err)
	assert.Equal(t, l.Name, lesson.Name)
	assert.Equal(t, l.Description, lesson.Description)
	assert.Equal(t, 2, lesson.ExerciseCount)
}

func TestFetchExercises(t *testing.T) {
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

	// to check of exercise without results will also be fetched
	exercise2 := &Exercise{LessonId: exercise.LessonId}
	createExercise(db, exercise2)

	// belongs to exercise, to be included
	result1 := &Result{ExerciseId: exercise.Id}
	createResult(db, result1)

	// does not belong to exercise, to be excluded
	result2 := &Result{}
	createResult(db, result2)

	res, err := r.FetchExercises(models.Lesson{Id: exercise.LessonId})

	assert.NoError(t, err)
	assert.IsType(t, models.Exercises{}, res)
	assert.Len(t, res, 2)
	assert.Equal(t, exercise.Id, res[1].Id)
	assert.Equal(t, exercise.Question, res[1].Question)
	assert.Equal(t, exercise.Answer, res[1].Answer)
	assert.Len(t, res[1].Results, 1)
	assert.Empty(t, res[0].Results)
	assert.Equal(t, result1.Id, res[1].Results[0].Id)
	assert.Equal(t, result1.Type, res[1].Results[0].Type)
	assert.Equal(t, result1.CreatedAt, res[1].Results[0].CreatedAt)
}
