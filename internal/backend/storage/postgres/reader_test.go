package postgres

import (
	"context"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage/entities"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAllExercises(t *testing.T) {
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

	exercise := &entities.Exercise{}
	createExercise(db, exercise)

	res, err := r.AllExercises()

	assert.NoError(t, err)
	assert.IsType(t, models.Exercises{}, res)
	assert.Len(t, res, 1)
	assert.Equal(t, exercise.Id, res[0].Id)
	assert.Equal(t, exercise.Question, res[0].Question)
	assert.Equal(t, exercise.Answer, res[0].Answer)
	assert.Equal(t, 0, res[0].BadAnswers)
	assert.Equal(t, 0, res[0].GoodAnswers)
}

func TestAllLessons(t *testing.T) {
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

	lesson := &entities.Lesson{}
	createLesson(db, lesson)

	res, err := r.AllLessons()

	assert.NoError(t, err)
	assert.IsType(t, models.Lessons{}, res)
	assert.Len(t, res, 1)
	assert.Equal(t, lesson.Id, res[0].Id)
	assert.Equal(t, lesson.Name, res[0].Name)
}

func TestRandomExercise(t *testing.T) {
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

	exercise := &entities.Exercise{}
	createExercise(db, exercise)

	res, err := r.RandomExercise()

	assert.NoError(t, err)
	assert.IsType(t, models.Exercise{}, res)
	assert.Equal(t, exercise.Id, res.Id)
	assert.Equal(t, exercise.Question, res.Question)
	assert.Equal(t, exercise.Answer, res.Answer)
	assert.Equal(t, 0, res.BadAnswers)
	assert.Equal(t, 0, res.GoodAnswers)
}
