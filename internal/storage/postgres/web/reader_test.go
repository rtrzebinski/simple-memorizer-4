package web

import (
	"context"
	"time"

	"github.com/guregu/null/v5"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
)

func (s *PostgresSuite) TestReader_FetchLessons() {
	db := s.db

	r := NewReader(db)

	user := &postgres.User{
		Id: 1,
	}
	postgres.CreateUser(db, user)

	lesson := &postgres.Lesson{
		UserID: 1,
	}
	postgres.CreateLesson(db, lesson)

	postgres.CreateExercise(db, &postgres.Exercise{LessonId: lesson.Id})

	res, err := r.FetchLessons(context.Background(), "1")

	assert.NoError(s.T(), err)
	assert.IsType(s.T(), backend.Lessons{}, res)
	assert.Len(s.T(), res, 1)
	assert.Equal(s.T(), lesson.Id, res[0].Id)
	assert.Equal(s.T(), lesson.Name, res[0].Name)
	assert.Equal(s.T(), lesson.Description, res[0].Description)
	assert.Equal(s.T(), 1, res[0].ExerciseCount)
}

func (s *PostgresSuite) TestReader_FetchLessons_otherUser() {
	db := s.db

	r := NewReader(db)

	user := &postgres.User{
		Id: 1,
	}
	postgres.CreateUser(db, user)

	lesson := &postgres.Lesson{
		UserID: 1,
	}
	postgres.CreateLesson(db, lesson)

	postgres.CreateExercise(db, &postgres.Exercise{LessonId: lesson.Id})

	res, err := r.FetchLessons(context.Background(), "2")

	assert.NoError(s.T(), err)
	assert.IsType(s.T(), backend.Lessons{}, res)
	assert.Len(s.T(), res, 0)
}

func (s *PostgresSuite) TestReader_HydrateLesson() {
	db := s.db

	r := NewReader(db)

	l := &postgres.Lesson{}
	postgres.CreateLesson(db, l)

	lesson := &backend.Lesson{
		Id: l.Id,
	}

	err := r.HydrateLesson(context.Background(), lesson, "userID")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), l.Name, lesson.Name)
	assert.Equal(s.T(), l.Description, lesson.Description)
	assert.Equal(s.T(), 0, lesson.ExerciseCount)

	postgres.CreateExercise(db, &postgres.Exercise{LessonId: l.Id})
	postgres.CreateExercise(db, &postgres.Exercise{LessonId: l.Id})

	err = r.HydrateLesson(context.Background(), lesson, "userID")

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), l.Name, lesson.Name)
	assert.Equal(s.T(), l.Description, lesson.Description)
	assert.Equal(s.T(), 2, lesson.ExerciseCount)
}

func (s *PostgresSuite) TestReader_FetchExercises() {
	db := s.db

	r := NewReader(db)

	exercise1 := &postgres.Exercise{
		BadAnswers:               1,
		BadAnswersToday:          2,
		LatestBadAnswer:          null.TimeFrom(time.Now()),
		GoodAnswers:              3,
		GoodAnswersToday:         4,
		LatestGoodAnswer:         null.Time{},
		LatestGoodAnswerWasToday: true,
	}
	postgres.CreateExercise(db, exercise1)

	// to check of exercise without results will also be fetched
	exercise2 := &postgres.Exercise{LessonId: exercise1.LessonId}
	postgres.CreateExercise(db, exercise2)

	res, err := r.FetchExercises(context.Background(), backend.Lesson{Id: exercise1.LessonId}, "userID")

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
