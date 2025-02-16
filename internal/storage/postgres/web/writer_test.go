package web

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
)

func (s *PostgresSuite) TestWriter_UpsertLesson_createNew() {
	db := s.db

	w := NewWriter(db)

	lesson := backend.Lesson{
		Name:        "name",
		Description: "description",
	}

	err := w.UpsertLesson(context.Background(), &lesson)
	assert.NoError(s.T(), err)

	stored := postgres.FetchLatestLesson(db)

	assert.Equal(s.T(), lesson.Name, stored.Name)
	assert.Equal(s.T(), lesson.Description, stored.Description)
	assert.Equal(s.T(), lesson.Id, stored.Id)
}

func (s *PostgresSuite) TestWriter_UpsertLesson_updateExisting() {
	db := s.db

	w := NewWriter(db)

	lesson := &postgres.Lesson{}
	postgres.CreateLesson(db, lesson)

	err := w.UpsertLesson(context.Background(), &backend.Lesson{
		Id:          1,
		Name:        "newName",
		Description: "newDescription",
	})
	assert.NoError(s.T(), err)

	stored := postgres.FetchLatestLesson(db)

	assert.Equal(s.T(), "newName", stored.Name)
	assert.Equal(s.T(), "newDescription", stored.Description)
}

func (s *PostgresSuite) TestWriter_DeleteLesson() {
	db := s.db

	w := NewWriter(db)

	postgres.CreateLesson(db, &postgres.Lesson{})
	stored := postgres.FetchLatestLesson(db)

	postgres.CreateLesson(db, &postgres.Lesson{
		Name: "another",
	})
	another := postgres.FetchLatestLesson(db)

	err := w.DeleteLesson(context.Background(), backend.Lesson{Id: stored.Id})
	assert.NoError(s.T(), err)

	assert.Nil(s.T(), postgres.FindLessonById(db, stored.Id))
	assert.Equal(s.T(), "another", postgres.FindLessonById(db, another.Id).Name)
}

func (s *PostgresSuite) TestWriter_UpsertExercise_createNew() {
	db := s.db

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
	assert.NoError(s.T(), err)

	stored := postgres.FetchLatestExercise(db)

	assert.Equal(s.T(), exercise.Lesson.Id, stored.LessonId)
	assert.Equal(s.T(), exercise.Question, stored.Question)
	assert.Equal(s.T(), exercise.Answer, stored.Answer)
	assert.Equal(s.T(), exercise.Id, stored.Id)
}

func (s *PostgresSuite) TestWriter_UpsertExercise_updateExisting() {
	db := s.db

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
	assert.NoError(s.T(), err)

	stored := postgres.FetchLatestExercise(db)

	assert.Equal(s.T(), lesson.Id, stored.LessonId)
	assert.Equal(s.T(), "newQuestion", stored.Question)
	assert.Equal(s.T(), "newAnswer", stored.Answer)
}

func (s *PostgresSuite) TestWriter_StoreExercises() {
	db := s.db

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
	assert.NoError(s.T(), err)

	ex1 := postgres.FindExerciseById(db, 1)

	assert.Equal(s.T(), exercise1.Lesson.Id, ex1.LessonId)
	assert.Equal(s.T(), exercise1.Question, ex1.Question)
	assert.Equal(s.T(), exercise1.Answer, ex1.Answer)

	// ID of inserted exercise will be 3, not 2,
	// this is because 'ON CONFLICT (lesson_id, question) DO NOTHING',
	// is still increasing PK auto increment value, even if nothing is inserted
	ex2 := postgres.FindExerciseById(db, 3)
	assert.Equal(s.T(), exercise2.Lesson.Id, ex2.LessonId)
	assert.Equal(s.T(), exercise2.Question, ex2.Question)
	assert.Equal(s.T(), exercise2.Answer, ex2.Answer)
}

func (s *PostgresSuite) TestWriter_DeleteExercise() {
	db := s.db

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
	assert.NoError(s.T(), err)

	assert.Nil(s.T(), postgres.FindExerciseById(db, stored.Id))
	assert.Equal(s.T(), "another", postgres.FindExerciseById(db, another.Id).Question)
}
