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

	user := &user{
		Id: 1,
	}
	createUser(s.DB, user)

	lesson := backend.Lesson{
		Name:        "name",
		Description: "description",
	}

	err := s.writer.UpsertLesson(ctx, &lesson, "1")
	assert.NoError(s.T(), err)

	stored := fetchLatestLesson(s.DB)

	assert.Equal(s.T(), lesson.Name, stored.Name)
	assert.Equal(s.T(), lesson.Description, stored.Description)
	assert.Equal(s.T(), lesson.Id, stored.Id)
	assert.Equal(s.T(), user.Id, stored.UserID)
}

func (s *WebWriterSuite) TestWebWriter_UpsertLesson_updateExisting() {
	ctx := s.T().Context()

	lesson := &lesson{}
	createLesson(s.DB, lesson)

	err := s.writer.UpsertLesson(ctx, &backend.Lesson{
		Id:          1,
		Name:        "newName",
		Description: "newDescription",
	}, "userID")
	assert.NoError(s.T(), err)

	stored := fetchLatestLesson(s.DB)

	assert.Equal(s.T(), "newName", stored.Name)
	assert.Equal(s.T(), "newDescription", stored.Description)
}

func (s *WebWriterSuite) TestWebWriter_DeleteLesson() {
	ctx := s.T().Context()

	createLesson(s.DB, &lesson{})
	stored := fetchLatestLesson(s.DB)

	createLesson(s.DB, &lesson{
		Name: "another",
	})
	another := fetchLatestLesson(s.DB)

	err := s.writer.DeleteLesson(ctx, backend.Lesson{Id: stored.Id}, "userID")
	assert.NoError(s.T(), err)

	assert.Nil(s.T(), findLessonById(s.DB, stored.Id))
	assert.Equal(s.T(), "another", findLessonById(s.DB, another.Id).Name)
}

func (s *WebWriterSuite) TestWebWriter_UpsertExercise_createNew() {
	ctx := s.T().Context()

	lesson := &lesson{}
	createLesson(s.DB, lesson)

	exercise := backend.Exercise{
		Lesson: &backend.Lesson{
			Id: lesson.Id,
		},
		Question: "question",
		Answer:   "answer",
	}

	err := s.writer.UpsertExercise(ctx, &exercise, "userID")
	assert.NoError(s.T(), err)

	stored := fetchLatestExercise(s.DB)

	assert.Equal(s.T(), exercise.Lesson.Id, stored.LessonId)
	assert.Equal(s.T(), exercise.Question, stored.Question)
	assert.Equal(s.T(), exercise.Answer, stored.Answer)
	assert.Equal(s.T(), exercise.Id, stored.Id)
}

func (s *WebWriterSuite) TestWebWriter_UpsertExercise_updateExisting() {
	ctx := s.T().Context()

	lesson := lesson{}
	createLesson(s.DB, &lesson)

	exercise := exercise{LessonId: lesson.Id}
	createExercise(s.DB, &exercise)

	err := s.writer.UpsertExercise(ctx, &backend.Exercise{
		Id:       1,
		Question: "newQuestion",
		Answer:   "newAnswer",
	}, "userID")
	assert.NoError(s.T(), err)

	stored := fetchLatestExercise(s.DB)

	assert.Equal(s.T(), lesson.Id, stored.LessonId)
	assert.Equal(s.T(), "newQuestion", stored.Question)
	assert.Equal(s.T(), "newAnswer", stored.Answer)
}

func (s *WebWriterSuite) TestWebWriter_StoreExercises() {
	ctx := s.T().Context()

	lesson := &lesson{}
	createLesson(s.DB, lesson)

	// exercise1 existing
	exercise1 := backend.Exercise{
		Lesson: &backend.Lesson{
			Id: lesson.Id,
		},
		Question: "question1",
		Answer:   "answer1",
	}

	// store exercise 1 to db
	createExercise(s.DB, &exercise{
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

	ex1 := findExerciseById(s.DB, 1)

	assert.Equal(s.T(), exercise1.Lesson.Id, ex1.LessonId)
	assert.Equal(s.T(), exercise1.Question, ex1.Question)
	assert.Equal(s.T(), exercise1.Answer, ex1.Answer)

	// ID of inserted exercise will be 3, not 2,
	// this is because 'ON CONFLICT (lesson_id, question) DO NOTHING',
	// is still increasing PK auto increment value, even if nothing is inserted
	ex2 := findExerciseById(s.DB, 3)
	assert.Equal(s.T(), exercise2.Lesson.Id, ex2.LessonId)
	assert.Equal(s.T(), exercise2.Question, ex2.Question)
	assert.Equal(s.T(), exercise2.Answer, ex2.Answer)
}

func (s *WebWriterSuite) TestWebWriter_DeleteExercise() {
	ctx := s.T().Context()

	lesson := &lesson{}
	createLesson(s.DB, lesson)

	createExercise(s.DB, &exercise{
		LessonId: lesson.Id,
	})
	stored := fetchLatestExercise(s.DB)

	createExercise(s.DB, &exercise{
		Question: "another",
	})
	another := fetchLatestExercise(s.DB)

	err := s.writer.DeleteExercise(ctx, backend.Exercise{Id: stored.Id}, "userID")
	assert.NoError(s.T(), err)

	assert.Nil(s.T(), findExerciseById(s.DB, stored.Id))
	assert.Equal(s.T(), "another", findExerciseById(s.DB, another.Id).Question)
}
