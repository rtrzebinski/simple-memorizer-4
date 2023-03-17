package postgres

import (
	"context"
	"github.com/rtrzebinski/simple-memorizer-4/internal/server/storage/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncrementBadAnswers(t *testing.T) {
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
			w.IncrementBadAnswers(exercise.Id)

			exerciseResult := findExerciseResultByExerciseId(db, exercise.Id)

			assert.Equal(t, 1, exerciseResult.BadAnswers)
			assert.Equal(t, 0, exerciseResult.GoodAnswers)
		})

	t.Run(
		"existing exercise result", func(t *testing.T) {
			w.IncrementBadAnswers(exercise.Id)

			exerciseResult := findExerciseResultByExerciseId(db, exercise.Id)

			assert.Equal(t, 2, exerciseResult.BadAnswers)
			assert.Equal(t, 0, exerciseResult.GoodAnswers)
		})
}

func TestIncrementGoodAnswers(t *testing.T) {
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
			w.IncrementGoodAnswers(exercise.Id)

			exerciseResult := findExerciseResultByExerciseId(db, exercise.Id)

			assert.Equal(t, 1, exerciseResult.GoodAnswers)
			assert.Equal(t, 0, exerciseResult.BadAnswers)
		})

	t.Run(
		"existing exercise result", func(t *testing.T) {
			w.IncrementGoodAnswers(exercise.Id)

			exerciseResult := findExerciseResultByExerciseId(db, exercise.Id)

			assert.Equal(t, 2, exerciseResult.GoodAnswers)
			assert.Equal(t, 0, exerciseResult.BadAnswers)
		})
}
