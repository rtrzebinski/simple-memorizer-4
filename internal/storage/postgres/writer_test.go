package postgres

import (
	"context"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIncrementBadAnswers_notExistingExerciseResult(t *testing.T) {
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

	w.IncrementBadAnswers(exercise.Id)

	exerciseResult := findExerciseResultByExerciseId(db, exercise.Id)

	assert.Equal(t, 1, exerciseResult.BadAnswers)
	assert.Equal(t, 0, exerciseResult.GoodAnswers)
}
