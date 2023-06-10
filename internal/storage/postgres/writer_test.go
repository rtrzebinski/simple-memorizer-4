package postgres

import (
	"context"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStoreLesson_createNew(t *testing.T) {
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

	lesson := models.Lesson{
		Name: "name",
	}

	err = w.StoreLesson(&lesson)
	assert.NoError(t, err)

	stored := fetchLatestLesson(db)

	assert.Equal(t, lesson.Name, stored.Name)
	assert.Equal(t, lesson.Id, stored.Id)
}

func TestStoreLesson_updateExisting(t *testing.T) {
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

	lesson := &Lesson{}
	createLesson(db, lesson)

	err = w.StoreLesson(&models.Lesson{
		Id:   1,
		Name: "newName",
	})
	assert.NoError(t, err)

	stored := fetchLatestLesson(db)

	assert.Equal(t, "newName", stored.Name)
}

func TestDeleteLesson(t *testing.T) {
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

	createLesson(db, &Lesson{})
	stored := fetchLatestLesson(db)

	createLesson(db, &Lesson{
		Name: "another",
	})
	another := fetchLatestLesson(db)

	err = w.DeleteLesson(models.Lesson{Id: stored.Id})
	assert.NoError(t, err)

	assert.Nil(t, findLessonById(db, stored.Id))
	assert.Equal(t, "another", findLessonById(db, another.Id).Name)
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

	lesson := &Lesson{}
	createLesson(db, lesson)

	exercise := models.Exercise{
		Lesson: &models.Lesson{
			Id: lesson.Id,
		},
		Question: "question",
		Answer:   "answer",
	}

	err = w.StoreExercise(exercise)
	assert.NoError(t, err)

	stored := fetchLatestExercise(db)

	assert.Equal(t, exercise.Lesson.Id, stored.LessonId)
	assert.Equal(t, exercise.Question, stored.Question)
	assert.Equal(t, exercise.Answer, stored.Answer)

	lesson = findLessonById(db, lesson.Id)
	assert.Equal(t, 1, lesson.ExerciseCount)
}

func TestStoreExercise_updateExisting(t *testing.T) {
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

	lesson := Lesson{}
	createLesson(db, &lesson)

	exercise := Exercise{LessonId: lesson.Id}
	createExercise(db, &exercise)

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

	lesson := &Lesson{
		ExerciseCount: 1,
	}
	createLesson(db, lesson)

	createExercise(db, &Exercise{
		LessonId: lesson.Id,
	})
	stored := fetchLatestExercise(db)

	createExercise(db, &Exercise{
		Question: "another",
	})
	another := fetchLatestExercise(db)

	err = w.DeleteExercise(models.Exercise{Id: stored.Id})
	assert.NoError(t, err)

	assert.Nil(t, findExerciseById(db, stored.Id))
	assert.Equal(t, "another", findExerciseById(db, another.Id).Question)

	lesson = findLessonById(db, stored.LessonId)
	assert.Equal(t, 0, lesson.ExerciseCount)
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

	exercise := &Exercise{}
	createExercise(db, exercise)

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

	exercise := &Exercise{}
	createExercise(db, exercise)

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
