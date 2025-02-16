package web

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
)

func (suite *PostgresSuite) TestWriter_UpsertLesson_createNew() {
	db := suite.db

	w := NewWriter(db)

	lesson := backend.Lesson{
		Name:        "name",
		Description: "description",
	}

	err := w.UpsertLesson(context.Background(), &lesson)
	assert.NoError(suite.T(), err)

	stored := postgres.FetchLatestLesson(db)

	assert.Equal(suite.T(), lesson.Name, stored.Name)
	assert.Equal(suite.T(), lesson.Description, stored.Description)
	assert.Equal(suite.T(), lesson.Id, stored.Id)
}

func (suite *PostgresSuite) TestWriter_UpsertLesson_updateExisting() {
	db := suite.db

	w := NewWriter(db)

	lesson := &postgres.Lesson{}
	postgres.CreateLesson(db, lesson)

	err := w.UpsertLesson(context.Background(), &backend.Lesson{
		Id:          1,
		Name:        "newName",
		Description: "newDescription",
	})
	assert.NoError(suite.T(), err)

	stored := postgres.FetchLatestLesson(db)

	assert.Equal(suite.T(), "newName", stored.Name)
	assert.Equal(suite.T(), "newDescription", stored.Description)
}

func (suite *PostgresSuite) TestWriter_DeleteLesson() {
	db := suite.db

	w := NewWriter(db)

	postgres.CreateLesson(db, &postgres.Lesson{})
	stored := postgres.FetchLatestLesson(db)

	postgres.CreateLesson(db, &postgres.Lesson{
		Name: "another",
	})
	another := postgres.FetchLatestLesson(db)

	err := w.DeleteLesson(context.Background(), backend.Lesson{Id: stored.Id})
	assert.NoError(suite.T(), err)

	assert.Nil(suite.T(), postgres.FindLessonById(db, stored.Id))
	assert.Equal(suite.T(), "another", postgres.FindLessonById(db, another.Id).Name)
}

func (suite *PostgresSuite) TestWriter_UpsertExercise_createNew() {
	db := suite.db

	w := NewWriter(db)

	lesson := &postgres.Lesson{}
	postgres.CreateLesson(db, lesson)

	exercise := backend.Exercise{
		Lesson: &backend.Lesson{
			Id: lesson.Id,
		},
		Question: "question",
		Answer:   "answer",
	}

	err := w.UpsertExercise(context.Background(), &exercise)
	assert.NoError(suite.T(), err)

	stored := postgres.FetchLatestExercise(db)

	assert.Equal(suite.T(), exercise.Lesson.Id, stored.LessonId)
	assert.Equal(suite.T(), exercise.Question, stored.Question)
	assert.Equal(suite.T(), exercise.Answer, stored.Answer)
	assert.Equal(suite.T(), exercise.Id, stored.Id)
}

func (suite *PostgresSuite) TestWriter_UpsertExercise_updateExisting() {
	db := suite.db

	w := NewWriter(db)

	lesson := postgres.Lesson{}
	postgres.CreateLesson(db, &lesson)

	exercise := postgres.Exercise{LessonId: lesson.Id}
	postgres.CreateExercise(db, &exercise)

	err := w.UpsertExercise(context.Background(), &backend.Exercise{
		Id:       1,
		Question: "newQuestion",
		Answer:   "newAnswer",
	})
	assert.NoError(suite.T(), err)

	stored := postgres.FetchLatestExercise(db)

	assert.Equal(suite.T(), lesson.Id, stored.LessonId)
	assert.Equal(suite.T(), "newQuestion", stored.Question)
	assert.Equal(suite.T(), "newAnswer", stored.Answer)
}

func (suite *PostgresSuite) TestWriter_StoreExercises() {
	db := suite.db

	w := NewWriter(db)

	lesson := &postgres.Lesson{}
	postgres.CreateLesson(db, lesson)

	// exercise1 existing
	exercise1 := backend.Exercise{
		Lesson: &backend.Lesson{
			Id: lesson.Id,
		},
		Question: "question1",
		Answer:   "answer1",
	}

	// store exercise 1 to db
	postgres.CreateExercise(db, &postgres.Exercise{
		LessonId: exercise1.Lesson.Id,
		Question: exercise1.Question,
		Answer:   exercise1.Answer,
	})

	// exercise2 not existing
	exercise2 := backend.Exercise{
		Lesson: &backend.Lesson{
			Id: lesson.Id,
		},
		Question: "question2",
		Answer:   "answer2",
	}

	exercises := backend.Exercises{exercise1, exercise2}

	err := w.StoreExercises(context.Background(), exercises)
	assert.NoError(suite.T(), err)

	ex1 := postgres.FindExerciseById(db, 1)

	assert.Equal(suite.T(), exercise1.Lesson.Id, ex1.LessonId)
	assert.Equal(suite.T(), exercise1.Question, ex1.Question)
	assert.Equal(suite.T(), exercise1.Answer, ex1.Answer)

	// ID of inserted exercise will be 3, not 2,
	// this is because 'ON CONFLICT (lesson_id, question) DO NOTHING',
	// is still increasing PK auto increment value, even if nothing is inserted
	ex2 := postgres.FindExerciseById(db, 3)
	assert.Equal(suite.T(), exercise2.Lesson.Id, ex2.LessonId)
	assert.Equal(suite.T(), exercise2.Question, ex2.Question)
	assert.Equal(suite.T(), exercise2.Answer, ex2.Answer)
}

func (suite *PostgresSuite) TestWriter_DeleteExercise() {
	db := suite.db

	w := NewWriter(db)

	lesson := &postgres.Lesson{}
	postgres.CreateLesson(db, lesson)

	postgres.CreateExercise(db, &postgres.Exercise{
		LessonId: lesson.Id,
	})
	stored := postgres.FetchLatestExercise(db)

	postgres.CreateExercise(db, &postgres.Exercise{
		Question: "another",
	})
	another := postgres.FetchLatestExercise(db)

	err := w.DeleteExercise(context.Background(), backend.Exercise{Id: stored.Id})
	assert.NoError(suite.T(), err)

	assert.Nil(suite.T(), postgres.FindExerciseById(db, stored.Id))
	assert.Equal(suite.T(), "another", postgres.FindExerciseById(db, another.Id).Question)
}
