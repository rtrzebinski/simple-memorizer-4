package postgres

import (
	"context"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage/entities"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDeleteExercise(t *testing.T) {
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

	w := NewWriter(db)

	storeExercise(db, &entities.Exercise{
		Question: "question",
		Answer:   "answer",
	})
	stored := fetchLatestExercise(db)

	storeExercise(db, &entities.Exercise{
		Question: "another",
		Answer:   "another",
	})
	another := fetchLatestExercise(db)

	err = w.DeleteExercise(models.Exercise{Id: stored.Id})
	assert.NoError(t, err)

	assert.Nil(t, findExerciseById(db, stored.Id))
	assert.Equal(t, "another", findExerciseById(db, another.Id).Question)
}

func TestStoreExercise_createNew(t *testing.T) {
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

	w := NewWriter(db)

	exercise := models.Exercise{
		Question: "question",
		Answer:   "answer",
	}

	err = w.StoreExercise(exercise)
	assert.NoError(t, err)

	stored := fetchLatestExercise(db)

	assert.Equal(t, exercise.Question, stored.Question)
	assert.Equal(t, exercise.Answer, stored.Answer)
}

func TestStoreExercise_editExisting(t *testing.T) {
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

	w := NewWriter(db)

	exercise := &entities.Exercise{
		Question: "question",
		Answer:   "answer",
	}

	storeExercise(db, exercise)

	err = w.StoreExercise(models.Exercise{
		Id:       1,
		Question: "newQuestion",
		Answer:   "newAnswer",
	})
	assert.NoError(t, err)

	stored := fetchLatestExercise(db)

	assert.Equal(t, "newQuestion", stored.Question)
	assert.Equal(t, "newAnswer", stored.Answer)
}

func TestIncrementBadAnswers(t *testing.T) {
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

	w := NewWriter(db)

	exercise := &entities.Exercise{
		Question: "question",
		Answer:   "answer",
	}

	storeExercise(db, exercise)

	t.Run(
		"not existing exercise result", func(t *testing.T) {
			err := w.IncrementBadAnswers(models.Exercise{Id: exercise.Id})
			assert.NoError(t, err)

			exerciseResult := findExerciseResultByExerciseId(db, exercise.Id)

			assert.Equal(t, 1, exerciseResult.BadAnswers)
			assert.Equal(t, 0, exerciseResult.GoodAnswers)
		})

	t.Run(
		"existing exercise result", func(t *testing.T) {
			err := w.IncrementBadAnswers(models.Exercise{Id: exercise.Id})
			assert.NoError(t, err)

			exerciseResult := findExerciseResultByExerciseId(db, exercise.Id)

			assert.Equal(t, 2, exerciseResult.BadAnswers)
			assert.Equal(t, 0, exerciseResult.GoodAnswers)
		})
}

func TestIncrementGoodAnswers(t *testing.T) {
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

	w := NewWriter(db)

	exercise := &entities.Exercise{
		Question: "question",
		Answer:   "answer",
	}

	storeExercise(db, exercise)

	t.Run(
		"not existing exercise result", func(t *testing.T) {
			err := w.IncrementGoodAnswers(models.Exercise{Id: exercise.Id})
			assert.NoError(t, err)

			exerciseResult := findExerciseResultByExerciseId(db, exercise.Id)

			assert.Equal(t, 1, exerciseResult.GoodAnswers)
			assert.Equal(t, 0, exerciseResult.BadAnswers)
		})

	t.Run(
		"existing exercise result", func(t *testing.T) {
			err := w.IncrementGoodAnswers(models.Exercise{Id: exercise.Id})
			assert.NoError(t, err)

			exerciseResult := findExerciseResultByExerciseId(db, exercise.Id)

			assert.Equal(t, 2, exerciseResult.GoodAnswers)
			assert.Equal(t, 0, exerciseResult.BadAnswers)
		})
}
