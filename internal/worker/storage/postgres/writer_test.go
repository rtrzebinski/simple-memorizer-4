package postgres

import (
	"time"

	"github.com/rtrzebinski/simple-memorizer-4/internal/worker"
	"github.com/stretchr/testify/assert"
)

func (suite *PostgresSuite) TestWriter_StoreAnswer() {
	db := suite.db

	w := NewWriter(db)

	exercise := &Exercise{}
	createExercise(db, exercise)

	result := worker.Result{
		Type:       worker.Good,
		ExerciseId: exercise.Id,
	}

	err := w.StoreResult(&result)
	assert.NoError(suite.T(), err)

	stored := fetchLatestResult(db)

	assert.Equal(suite.T(), string(result.Type), stored.Type)
	assert.Equal(suite.T(), result.ExerciseId, stored.ExerciseId)
}

func (suite *PostgresSuite) TestWriter_UpdateExerciseProjection() {
	db := suite.db

	w := NewWriter(db)

	exercise := &Exercise{}
	createExercise(db, exercise)

	projection := worker.ResultsProjection{
		BadAnswers:               1,
		BadAnswersToday:          2,
		LatestBadAnswer:          time.Date(2021, time.March, 1, 12, 30, 0, 0, time.UTC),
		LatestBadAnswerWasToday:  true,
		GoodAnswers:              3,
		GoodAnswersToday:         4,
		LatestGoodAnswer:         time.Date(2022, time.March, 1, 12, 30, 0, 0, time.UTC),
		LatestGoodAnswerWasToday: true,
	}

	err := w.UpdateExerciseProjection(exercise.Id, projection)
	assert.NoError(suite.T(), err)

	stored := findExerciseById(db, exercise.Id)

	assert.Equal(suite.T(), projection.BadAnswers, stored.BadAnswers)
	assert.Equal(suite.T(), projection.BadAnswersToday, stored.BadAnswersToday)
	assert.True(suite.T(), projection.LatestBadAnswer.Equal(stored.LatestBadAnswer))
	assert.Equal(suite.T(), projection.LatestBadAnswerWasToday, stored.LatestBadAnswerWasToday)
	assert.Equal(suite.T(), projection.GoodAnswers, stored.GoodAnswers)
	assert.Equal(suite.T(), projection.GoodAnswersToday, stored.GoodAnswersToday)
	assert.True(suite.T(), projection.LatestGoodAnswer.Equal(stored.LatestGoodAnswer))
	assert.Equal(suite.T(), projection.LatestGoodAnswerWasToday, stored.LatestGoodAnswerWasToday)
}
