package postgres

import (
	"testing"
	"time"

	"github.com/guregu/null/v5"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WebReaderSuite struct {
	Suite
	reader *WebReader
}

func TestWebReader(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	suite.Run(t, new(WebReaderSuite))
}

func (s *WebReaderSuite) SetupSuite() {
	s.Suite.SetupSuite()
	s.reader = NewWebReader(s.DB)
}

func (s *WebReaderSuite) TestWebReader_FetchLessons() {
	ctx := s.T().Context()

	userID := randomString()

	lesson := &lesson{
		UserID: userID,
	}
	createLesson(s.DB, lesson)

	createExercise(s.DB, &exercise{LessonId: lesson.Id})

	res, err := s.reader.FetchLessons(ctx, userID)

	assert.NoError(s.T(), err)
	assert.IsType(s.T(), backend.Lessons{}, res)
	assert.Len(s.T(), res, 1)
	assert.Equal(s.T(), lesson.Id, res[0].Id)
	assert.Equal(s.T(), lesson.Name, res[0].Name)
	assert.Equal(s.T(), lesson.Description, res[0].Description)
	assert.Equal(s.T(), 1, res[0].ExerciseCount)
}

func (s *WebReaderSuite) TestWebReader_FetchLessons_otherUser() {
	ctx := s.T().Context()

	lesson := &lesson{}
	createLesson(s.DB, lesson)

	createExercise(s.DB, &exercise{LessonId: lesson.Id})

	res, err := s.reader.FetchLessons(ctx, "2")

	assert.NoError(s.T(), err)
	assert.IsType(s.T(), backend.Lessons{}, res)
	assert.Len(s.T(), res, 0)
}

func (s *WebReaderSuite) TestWebReader_HydrateLesson() {
	ctx := s.T().Context()

	l := &lesson{}
	createLesson(s.DB, l)

	lesson := &backend.Lesson{
		Id: l.Id,
	}

	err := s.reader.HydrateLesson(ctx, lesson, "userID")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), l.Name, lesson.Name)
	assert.Equal(s.T(), l.Description, lesson.Description)
	assert.Equal(s.T(), 0, lesson.ExerciseCount)

	createExercise(s.DB, &exercise{LessonId: l.Id})
	createExercise(s.DB, &exercise{LessonId: l.Id})

	err = s.reader.HydrateLesson(ctx, lesson, "userID")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), l.Name, lesson.Name)
	assert.Equal(s.T(), l.Description, lesson.Description)
	assert.Equal(s.T(), 2, lesson.ExerciseCount)
}

func (s *WebReaderSuite) TestWebReader_FetchExercises() {
	ctx := s.T().Context()

	exercise1 := &exercise{
		BadAnswers:               1,
		BadAnswersToday:          2,
		LatestBadAnswer:          null.TimeFrom(time.Now()),
		GoodAnswers:              3,
		GoodAnswersToday:         4,
		LatestGoodAnswer:         null.Time{},
		LatestGoodAnswerWasToday: true,
	}
	createExercise(s.DB, exercise1)

	// to check of exercise without results will also be fetched
	exercise2 := &exercise{LessonId: exercise1.LessonId}
	createExercise(s.DB, exercise2)

	oldestExerciseID := 1

	res, err := s.reader.FetchExercises(ctx, backend.Lesson{Id: exercise1.LessonId}, oldestExerciseID, "userID")

	assert.NoError(s.T(), err)
	assert.IsType(s.T(), backend.Exercises{}, res)
	assert.Len(s.T(), res, 2)
	assert.Equal(s.T(), exercise1.Id, res[1].Id)
	assert.Equal(s.T(), exercise1.Question, res[1].Question)
	assert.Equal(s.T(), exercise1.Answer, res[1].Answer)
	assert.Equal(s.T(), exercise1.BadAnswers, res[1].BadAnswers)
	assert.Equal(s.T(), exercise1.BadAnswersToday, res[1].BadAnswersToday)
	assert.Equal(s.T(), exercise1.LatestBadAnswer.Time.Local().Format("Mon Jan 2 15:04:05"), res[1].LatestBadAnswer.Time.Local().Format("Mon Jan 2 15:04:05"))
	assert.Equal(s.T(), exercise1.GoodAnswers, res[1].GoodAnswers)
	assert.Equal(s.T(), exercise1.GoodAnswersToday, res[1].GoodAnswersToday)
	assert.Equal(s.T(), exercise1.LatestGoodAnswer, res[1].LatestGoodAnswer)
	assert.Equal(s.T(), exercise1.LatestGoodAnswerWasToday, res[1].LatestGoodAnswerWasToday)
	assert.Equal(s.T(), exercise2.Id, res[0].Id)
	assert.Equal(s.T(), exercise2.Question, res[0].Question)
	assert.Equal(s.T(), exercise2.Answer, res[0].Answer)
	assert.Equal(s.T(), 0, res[0].BadAnswers)
	assert.Equal(s.T(), 0, res[0].BadAnswersToday)
	assert.Equal(s.T(), null.Time{}, res[0].LatestBadAnswer)
	assert.Equal(s.T(), 0, res[0].GoodAnswers)
	assert.Equal(s.T(), 0, res[0].GoodAnswersToday)
	assert.Equal(s.T(), null.Time{}, res[0].LatestGoodAnswer)
	assert.Equal(s.T(), false, res[0].LatestGoodAnswerWasToday)
}

func (s *WebReaderSuite) TestWebReader_FetchExercises_oldestExerciseID() {
	ctx := s.T().Context()

	exercise1 := &exercise{
		BadAnswers:               1,
		BadAnswersToday:          2,
		LatestBadAnswer:          null.TimeFrom(time.Now()),
		GoodAnswers:              3,
		GoodAnswersToday:         4,
		LatestGoodAnswer:         null.Time{},
		LatestGoodAnswerWasToday: true,
	}
	createExercise(s.DB, exercise1)

	// to check of exercise without results will also be fetched
	exercise2 := &exercise{LessonId: exercise1.LessonId}
	createExercise(s.DB, exercise2)

	oldestExerciseID := 2

	res, err := s.reader.FetchExercises(ctx, backend.Lesson{Id: exercise1.LessonId}, oldestExerciseID, "userID")

	assert.NoError(s.T(), err)
	assert.IsType(s.T(), backend.Exercises{}, res)
	assert.Len(s.T(), res, 1)
	assert.Equal(s.T(), exercise2.Id, res[0].Id)
}
