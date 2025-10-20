package postgres

import (
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WebWriterSuite struct {
	PostgresSuite
	writer *WebWriter
}

func TestWebWriter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	suite.Run(t, new(WebWriterSuite))
}

func (s *WebWriterSuite) SetupSuite() {
	s.PostgresSuite.SetupSuite()
	s.writer = NewWebWriter(s.DB)
}

func (s *WebWriterSuite) TestWebWriter_UpsertLesson_createNew() {
	ctx := s.T().Context()

	user := &User{
		Id: 1,
	}
	CreateUser(s.DB, user)

	lesson := backend.Lesson{
		Name:        "name",
		Description: "description",
	}

	err := s.writer.UpsertLesson(ctx, &lesson, "1")
	assert.NoError(s.T(), err)

	stored := FetchLatestLesson(s.DB)

	assert.Equal(s.T(), lesson.Name, stored.Name)
	assert.Equal(s.T(), lesson.Description, stored.Description)
	assert.Equal(s.T(), lesson.Id, stored.Id)
	assert.Equal(s.T(), user.Id, stored.UserID)
}

func (s *WebWriterSuite) TestWebWriter_UpsertLesson_updateExisting() {
	ctx := s.T().Context()

	lesson := &Lesson{}
	CreateLesson(s.DB, lesson)

	err := s.writer.UpsertLesson(ctx, &backend.Lesson{
		Id:          1,
		Name:        "newName",
		Description: "newDescription",
	}, "userID")
	assert.NoError(s.T(), err)

	stored := FetchLatestLesson(s.DB)

	assert.Equal(s.T(), "newName", stored.Name)
	assert.Equal(s.T(), "newDescription", stored.Description)
}

func (s *WebWriterSuite) TestWebWriter_DeleteLesson() {
	ctx := s.T().Context()

	CreateLesson(s.DB, &Lesson{})
	stored := FetchLatestLesson(s.DB)

	CreateLesson(s.DB, &Lesson{
		Name: "another",
	})
	another := FetchLatestLesson(s.DB)

	err := s.writer.DeleteLesson(ctx, backend.Lesson{Id: stored.Id}, "userID")
	assert.NoError(s.T(), err)

	assert.Nil(s.T(), FindLessonById(s.DB, stored.Id))
	assert.Equal(s.T(), "another", FindLessonById(s.DB, another.Id).Name)
}

func (s *WebWriterSuite) TestWebWriter_UpsertExercise_createNew() {
	ctx := s.T().Context()

	lesson := &Lesson{}
	CreateLesson(s.DB, lesson)

	exercise := backend.Exercise{
		Lesson: &backend.Lesson{
			Id: lesson.Id,
		},
		Question: "question",
		Answer:   "answer",
	}

	err := s.writer.UpsertExercise(ctx, &exercise, "userID")
	assert.NoError(s.T(), err)

	stored := FetchLatestExercise(s.DB)

	assert.Equal(s.T(), exercise.Lesson.Id, stored.LessonId)
	assert.Equal(s.T(), exercise.Question, stored.Question)
	assert.Equal(s.T(), exercise.Answer, stored.Answer)
	assert.Equal(s.T(), exercise.Id, stored.Id)
}

func (s *WebWriterSuite) TestWebWriter_UpsertExercise_updateExisting() {
	ctx := s.T().Context()

	lesson := Lesson{}
	CreateLesson(s.DB, &lesson)

	exercise := Exercise{LessonId: lesson.Id}
	CreateExercise(s.DB, &exercise)

	err := s.writer.UpsertExercise(ctx, &backend.Exercise{
		Id:       1,
		Question: "newQuestion",
		Answer:   "newAnswer",
	}, "userID")
	assert.NoError(s.T(), err)

	stored := FetchLatestExercise(s.DB)

	assert.Equal(s.T(), lesson.Id, stored.LessonId)
	assert.Equal(s.T(), "newQuestion", stored.Question)
	assert.Equal(s.T(), "newAnswer", stored.Answer)
}

func (s *WebWriterSuite) TestWebWriter_StoreExercises() {
	ctx := s.T().Context()

	lesson := &Lesson{}
	CreateLesson(s.DB, lesson)

	// exercise1 existing
	exercise1 := backend.Exercise{
		Lesson: &backend.Lesson{
			Id: lesson.Id,
		},
		Question: "question1",
		Answer:   "answer1",
	}

	// store exercise 1 to db
	CreateExercise(s.DB, &Exercise{
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

	err := s.writer.StoreExercises(ctx, exercises, "userID")
	assert.NoError(s.T(), err)

	ex1 := FindExerciseById(s.DB, 1)

	assert.Equal(s.T(), exercise1.Lesson.Id, ex1.LessonId)
	assert.Equal(s.T(), exercise1.Question, ex1.Question)
	assert.Equal(s.T(), exercise1.Answer, ex1.Answer)

	// ID of inserted exercise will be 3, not 2,
	// this is because 'ON CONFLICT (lesson_id, question) DO NOTHING',
	// is still increasing PK auto increment value, even if nothing is inserted
	ex2 := FindExerciseById(s.DB, 3)
	assert.Equal(s.T(), exercise2.Lesson.Id, ex2.LessonId)
	assert.Equal(s.T(), exercise2.Question, ex2.Question)
	assert.Equal(s.T(), exercise2.Answer, ex2.Answer)
}

func (s *WebWriterSuite) TestWebWriter_DeleteExercise() {
	ctx := s.T().Context()

	lesson := &Lesson{}
	CreateLesson(s.DB, lesson)

	CreateExercise(s.DB, &Exercise{
		LessonId: lesson.Id,
	})
	stored := FetchLatestExercise(s.DB)

	CreateExercise(s.DB, &Exercise{
		Question: "another",
	})
	another := FetchLatestExercise(s.DB)

	err := s.writer.DeleteExercise(ctx, backend.Exercise{Id: stored.Id}, "userID")
	assert.NoError(s.T(), err)

	assert.Nil(s.T(), FindExerciseById(s.DB, stored.Id))
	assert.Equal(s.T(), "another", FindExerciseById(s.DB, another.Id).Question)
}
