package worker

import (
	"context"
	"time"

	"github.com/guregu/null/v5"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/worker"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
)

func (suite *PostgresSuite) TestWriter_StoreAnswer() {
	ctx := context.Background()

	db := suite.db

	w := NewWriter(db)

	exercise := &postgres.Exercise{}
	postgres.CreateExercise(db, exercise)

	result := worker.Result{
		Type:       worker.Good,
		ExerciseId: exercise.Id,
	}

	err := w.StoreResult(ctx, result)
	assert.NoError(suite.T(), err)

	stored := postgres.FetchLatestResult(db)

	assert.Equal(suite.T(), string(result.Type), stored.Type)
	assert.Equal(suite.T(), result.ExerciseId, stored.ExerciseId)
}

func (suite *PostgresSuite) TestWriter_UpdateExerciseProjection_allProjections() {
	ctx := context.Background()

	db := suite.db

	w := NewWriter(db)

	exercise := &postgres.Exercise{}
	postgres.CreateExercise(db, exercise)

	projection := worker.ResultsProjection{
		BadAnswers:               1,
		BadAnswersToday:          2,
		LatestBadAnswer:          null.TimeFrom(time.Date(2021, time.March, 1, 12, 30, 0, 0, time.UTC)),
		LatestBadAnswerWasToday:  true,
		GoodAnswers:              3,
		GoodAnswersToday:         4,
		LatestGoodAnswer:         null.TimeFrom(time.Date(2022, time.March, 1, 12, 30, 0, 0, time.UTC)),
		LatestGoodAnswerWasToday: true,
	}

	err := w.UpdateExerciseProjection(ctx, exercise.Id, projection)
	assert.NoError(suite.T(), err)

	stored := postgres.FindExerciseById(db, exercise.Id)

	assert.Equal(suite.T(), projection.BadAnswers, stored.BadAnswers)
	assert.Equal(suite.T(), projection.BadAnswersToday, stored.BadAnswersToday)
	assert.True(suite.T(), projection.LatestBadAnswer.Equal(stored.LatestBadAnswer))
	assert.Equal(suite.T(), projection.LatestBadAnswerWasToday, stored.LatestBadAnswerWasToday)
	assert.Equal(suite.T(), projection.GoodAnswers, stored.GoodAnswers)
	assert.Equal(suite.T(), projection.GoodAnswersToday, stored.GoodAnswersToday)
	assert.True(suite.T(), projection.LatestGoodAnswer.Equal(stored.LatestGoodAnswer))
	assert.Equal(suite.T(), projection.LatestGoodAnswerWasToday, stored.LatestGoodAnswerWasToday)
}

func (suite *PostgresSuite) TestWriter_UpdateExerciseProjection_badOnly() {
	ctx := context.Background()

	db := suite.db

	w := NewWriter(db)

	exercise := &postgres.Exercise{}
	postgres.CreateExercise(db, exercise)

	projection := worker.ResultsProjection{
		BadAnswers:              1,
		BadAnswersToday:         2,
		LatestBadAnswer:         null.TimeFrom(time.Date(2021, time.March, 1, 12, 30, 0, 0, time.UTC)),
		LatestBadAnswerWasToday: true,
	}

	err := w.UpdateExerciseProjection(ctx, exercise.Id, projection)
	assert.NoError(suite.T(), err)

	stored := postgres.FindExerciseById(db, exercise.Id)

	assert.Equal(suite.T(), projection.BadAnswers, stored.BadAnswers)
	assert.Equal(suite.T(), projection.BadAnswersToday, stored.BadAnswersToday)
	assert.True(suite.T(), projection.LatestBadAnswer.Equal(stored.LatestBadAnswer))
	assert.Equal(suite.T(), projection.LatestBadAnswerWasToday, stored.LatestBadAnswerWasToday)
	assert.Equal(suite.T(), projection.GoodAnswers, stored.GoodAnswers)
	assert.Equal(suite.T(), projection.GoodAnswersToday, stored.GoodAnswersToday)
	assert.True(suite.T(), projection.LatestGoodAnswer.Equal(stored.LatestGoodAnswer))
	assert.Equal(suite.T(), projection.LatestGoodAnswerWasToday, stored.LatestGoodAnswerWasToday)
}

func (suite *PostgresSuite) TestWriter_UpdateExerciseProjection_goodOnly() {
	ctx := context.Background()

	db := suite.db

	w := NewWriter(db)

	exercise := &postgres.Exercise{}
	postgres.CreateExercise(db, exercise)

	projection := worker.ResultsProjection{
		GoodAnswers:              3,
		GoodAnswersToday:         4,
		LatestGoodAnswer:         null.TimeFrom(time.Date(2022, time.March, 1, 12, 30, 0, 0, time.UTC)),
		LatestGoodAnswerWasToday: true,
	}

	err := w.UpdateExerciseProjection(ctx, exercise.Id, projection)
	assert.NoError(suite.T(), err)

	stored := postgres.FindExerciseById(db, exercise.Id)

	assert.Equal(suite.T(), projection.BadAnswers, stored.BadAnswers)
	assert.Equal(suite.T(), projection.BadAnswersToday, stored.BadAnswersToday)
	assert.True(suite.T(), projection.LatestBadAnswer.Equal(stored.LatestBadAnswer))
	assert.Equal(suite.T(), projection.LatestBadAnswerWasToday, stored.LatestBadAnswerWasToday)
	assert.Equal(suite.T(), projection.GoodAnswers, stored.GoodAnswers)
	assert.Equal(suite.T(), projection.GoodAnswersToday, stored.GoodAnswersToday)
	assert.True(suite.T(), projection.LatestGoodAnswer.Equal(stored.LatestGoodAnswer))
	assert.Equal(suite.T(), projection.LatestGoodAnswerWasToday, stored.LatestGoodAnswerWasToday)
}
