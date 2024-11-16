package postgres

import (
	"github.com/guregu/null/v5"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/stretchr/testify/assert"
	"time"
)

func (suite *PostgresSuite) TestReader_FetchLessons() {
	db := suite.db

	r := NewReader(db)

	lesson := &Lesson{}
	createLesson(db, lesson)

	createExercise(db, &Exercise{LessonId: lesson.Id})

	res, err := r.FetchLessons()

	assert.NoError(suite.T(), err)
	assert.IsType(suite.T(), backend.Lessons{}, res)
	assert.Len(suite.T(), res, 1)
	assert.Equal(suite.T(), lesson.Id, res[0].Id)
	assert.Equal(suite.T(), lesson.Name, res[0].Name)
	assert.Equal(suite.T(), lesson.Description, res[0].Description)
	assert.Equal(suite.T(), 1, res[0].ExerciseCount)
}

func (suite *PostgresSuite) TestReader_HydrateLesson() {
	db := suite.db

	r := NewReader(db)

	l := &Lesson{}
	createLesson(db, l)

	lesson := &backend.Lesson{
		Id: l.Id,
	}

	err := r.HydrateLesson(lesson)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), l.Name, lesson.Name)
	assert.Equal(suite.T(), l.Description, lesson.Description)
	assert.Equal(suite.T(), 0, lesson.ExerciseCount)

	createExercise(db, &Exercise{LessonId: l.Id})
	createExercise(db, &Exercise{LessonId: l.Id})

	err = r.HydrateLesson(lesson)

	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), l.Name, lesson.Name)
	assert.Equal(suite.T(), l.Description, lesson.Description)
	assert.Equal(suite.T(), 2, lesson.ExerciseCount)
}

func (suite *PostgresSuite) TestReader_FetchExercises() {
	db := suite.db

	r := NewReader(db)

	exercise1 := &Exercise{
		BadAnswers:               1,
		BadAnswersToday:          2,
		LatestBadAnswer:          null.TimeFrom(time.Now()),
		GoodAnswers:              3,
		GoodAnswersToday:         4,
		LatestGoodAnswer:         null.Time{},
		LatestGoodAnswerWasToday: true,
	}
	createExercise(db, exercise1)

	// to check of exercise without results will also be fetched
	exercise2 := &Exercise{LessonId: exercise1.LessonId}
	createExercise(db, exercise2)

	res, err := r.FetchExercises(backend.Lesson{Id: exercise1.LessonId})

	assert.NoError(suite.T(), err)
	assert.IsType(suite.T(), backend.Exercises{}, res)
	assert.Len(suite.T(), res, 2)
	assert.Equal(suite.T(), exercise1.Id, res[1].Id)
	assert.Equal(suite.T(), exercise1.Question, res[1].Question)
	assert.Equal(suite.T(), exercise1.Answer, res[1].Answer)
	assert.Equal(suite.T(), exercise1.BadAnswers, res[1].BadAnswers)
	assert.Equal(suite.T(), exercise1.BadAnswersToday, res[1].BadAnswersToday)
	assert.Equal(suite.T(), exercise1.LatestBadAnswer.Time.Local().Format("Mon Jan 2 15:04:05"), res[1].LatestBadAnswer.Time.Local().Format("Mon Jan 2 15:04:05"))
	assert.Equal(suite.T(), exercise1.GoodAnswers, res[1].GoodAnswers)
	assert.Equal(suite.T(), exercise1.GoodAnswersToday, res[1].GoodAnswersToday)
	assert.Equal(suite.T(), exercise1.LatestGoodAnswer, res[1].LatestGoodAnswer)
	assert.Equal(suite.T(), exercise1.LatestGoodAnswerWasToday, res[1].LatestGoodAnswerWasToday)
	assert.Equal(suite.T(), exercise2.Id, res[0].Id)
	assert.Equal(suite.T(), exercise2.Question, res[0].Question)
	assert.Equal(suite.T(), exercise2.Answer, res[0].Answer)
	assert.Equal(suite.T(), 0, res[0].BadAnswers)
	assert.Equal(suite.T(), 0, res[0].BadAnswersToday)
	assert.Equal(suite.T(), null.Time{}, res[0].LatestBadAnswer)
	assert.Equal(suite.T(), 0, res[0].GoodAnswers)
	assert.Equal(suite.T(), 0, res[0].GoodAnswersToday)
	assert.Equal(suite.T(), null.Time{}, res[0].LatestGoodAnswer)
	assert.Equal(suite.T(), false, res[0].LatestGoodAnswerWasToday)
}
