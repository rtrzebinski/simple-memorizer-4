package postgres

import (
	"context"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUpsertLesson_createNew(t *testing.T) {
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
		Name:        "name",
		Description: "description",
	}

	err = w.UpsertLesson(&lesson)
	assert.NoError(t, err)

	stored := fetchLatestLesson(db)

	assert.Equal(t, lesson.Name, stored.Name)
	assert.Equal(t, lesson.Description, stored.Description)
	assert.Equal(t, lesson.Id, stored.Id)
}

func TestUpsertLesson_updateExisting(t *testing.T) {
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

	err = w.UpsertLesson(&models.Lesson{
		Id:          1,
		Name:        "newName",
		Description: "newDescription",
	})
	assert.NoError(t, err)

	stored := fetchLatestLesson(db)

	assert.Equal(t, "newName", stored.Name)
	assert.Equal(t, "newDescription", stored.Description)
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

func TestUpsertExercise_createNew(t *testing.T) {
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

	err = w.UpsertExercise(&exercise)
	assert.NoError(t, err)

	stored := fetchLatestExercise(db)

	assert.Equal(t, exercise.Lesson.Id, stored.LessonId)
	assert.Equal(t, exercise.Question, stored.Question)
	assert.Equal(t, exercise.Answer, stored.Answer)
	assert.Equal(t, exercise.Id, stored.Id)
}

func TestUpsertExercise_updateExisting(t *testing.T) {
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

	err = w.UpsertExercise(&models.Exercise{
		Id:       1,
		Question: "newQuestion",
		Answer:   "newAnswer",
	})
	assert.NoError(t, err)

	stored := fetchLatestExercise(db)

	assert.Equal(t, "newQuestion", stored.Question)
	assert.Equal(t, "newAnswer", stored.Answer)
}

func TestStoreExercises(t *testing.T) {
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

	// exercise1 existing
	exercise1 := models.Exercise{
		Lesson: &models.Lesson{
			Id: lesson.Id,
		},
		Question: "question1",
		Answer:   "answer1",
	}

	// store exercise 1 to db
	createExercise(db, &Exercise{
		LessonId: exercise1.Lesson.Id,
		Question: exercise1.Question,
		Answer:   exercise1.Answer,
	})

	// exercise2 not existing
	exercise2 := models.Exercise{
		Lesson: &models.Lesson{
			Id: lesson.Id,
		},
		Question: "question2",
		Answer:   "answer2",
	}

	exercises := models.Exercises{exercise1, exercise2}

	err = w.StoreExercises(exercises)
	assert.NoError(t, err)

	ex1 := findExerciseById(db, 1)
	assert.Equal(t, exercise1.Lesson.Id, ex1.LessonId)
	assert.Equal(t, exercise1.Question, ex1.Question)
	assert.Equal(t, exercise1.Answer, ex1.Answer)

	// ID of inserted exercise will be 3, not 2,
	// this is because 'ON CONFLICT (lesson_id, question) DO NOTHING',
	// is still increasing PK auto increment value, even if nothing is inserted
	ex2 := findExerciseById(db, 3)
	assert.Equal(t, exercise2.Lesson.Id, ex2.LessonId)
	assert.Equal(t, exercise2.Question, ex2.Question)
	assert.Equal(t, exercise2.Answer, ex2.Answer)
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

func TestStoreAnswer(t *testing.T) {
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

	answer := models.Result{
		Type: models.Good,
		Exercise: &models.Exercise{
			Id: exercise.Id,
		},
	}

	err = w.StoreResult(&answer)
	assert.NoError(t, err)

	stored := fetchLatestResult(db)

	assert.Equal(t, answer.Type, stored.Type)
	assert.Equal(t, answer.Exercise.Id, stored.ExerciseId)
}
