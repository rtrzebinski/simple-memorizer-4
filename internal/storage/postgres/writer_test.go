package postgres

import (
	"context"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/entities"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRandomExercise(t *testing.T) {
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

	exercise := &entities.Exercise{
		Question: "question",
		Answer:   "answer",
	}

	storeExercise(db, exercise)

	res := r.RandomExercise()

	assert.IsType(t, models.Exercise{}, res)
	assert.Equal(t, exercise.Question, res.Question)
	assert.Equal(t, exercise.Answer, res.Answer)
}
