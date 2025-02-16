package worker

import (
	"context"
	"time"

	"github.com/guregu/null/v5"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/worker"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage/postgres"
	"github.com/stretchr/testify/assert"
)

func (s *PostgresSuite) TestWriter_StoreAnswer() {
	ctx := context.Background()

	db := s.db

	w := NewWriter(db)

	exercise := &postgres.Exercise{}
	postgres.CreateExercise(db, exercise)

	result := worker.Result{
		Type:       worker.Good,
		ExerciseId: exercise.Id,
	}

	err := w.StoreResult(ctx, result)
	assert.NoError(s.T(), err)

	stored := postgres.FetchLatestResult(db)

	assert.Equal(s.T(), string(result.Type), stored.Type)
	assert.Equal(s.T(), result.ExerciseId, stored.ExerciseId)
}

func (s *PostgresSuite) TestWriter_UpdateExerciseProjection_allProjections() {
	ctx := context.Background()

	db := s.db

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
	assert.NoError(s.T(), err)

	stored := postgres.FindExerciseById(db, exercise.Id)

	assert.Equal(s.T(), projection.BadAnswers, stored.BadAnswers)
	assert.Equal(s.T(), projection.BadAnswersToday, stored.BadAnswersToday)
	assert.True(s.T(), projection.LatestBadAnswer.Equal(stored.LatestBadAnswer))
	assert.Equal(s.T(), projection.LatestBadAnswerWasToday, stored.LatestBadAnswerWasToday)
	assert.Equal(s.T(), projection.GoodAnswers, stored.GoodAnswers)
	assert.Equal(s.T(), projection.GoodAnswersToday, stored.GoodAnswersToday)
	assert.True(s.T(), projection.LatestGoodAnswer.Equal(stored.LatestGoodAnswer))
	assert.Equal(s.T(), projection.LatestGoodAnswerWasToday, stored.LatestGoodAnswerWasToday)
}

func (s *PostgresSuite) TestWriter_UpdateExerciseProjection_badOnly() {
	ctx := context.Background()

	db := s.db

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
	assert.NoError(s.T(), err)

	stored := postgres.FindExerciseById(db, exercise.Id)

	assert.Equal(s.T(), projection.BadAnswers, stored.BadAnswers)
	assert.Equal(s.T(), projection.BadAnswersToday, stored.BadAnswersToday)
	assert.True(s.T(), projection.LatestBadAnswer.Equal(stored.LatestBadAnswer))
	assert.Equal(s.T(), projection.LatestBadAnswerWasToday, stored.LatestBadAnswerWasToday)
	assert.Equal(s.T(), projection.GoodAnswers, stored.GoodAnswers)
	assert.Equal(s.T(), projection.GoodAnswersToday, stored.GoodAnswersToday)
	assert.True(s.T(), projection.LatestGoodAnswer.Equal(stored.LatestGoodAnswer))
	assert.Equal(s.T(), projection.LatestGoodAnswerWasToday, stored.LatestGoodAnswerWasToday)
}

func (s *PostgresSuite) TestWriter_UpdateExerciseProjection_goodOnly() {
	ctx := context.Background()

	db := s.db

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
	assert.NoError(s.T(), err)

	stored := postgres.FindExerciseById(db, exercise.Id)

	assert.Equal(s.T(), projection.BadAnswers, stored.BadAnswers)
	assert.Equal(s.T(), projection.BadAnswersToday, stored.BadAnswersToday)
	assert.True(s.T(), projection.LatestBadAnswer.Equal(stored.LatestBadAnswer))
	assert.Equal(s.T(), projection.LatestBadAnswerWasToday, stored.LatestBadAnswerWasToday)
	assert.Equal(s.T(), projection.GoodAnswers, stored.GoodAnswers)
	assert.Equal(s.T(), projection.GoodAnswersToday, stored.GoodAnswersToday)
	assert.True(s.T(), projection.LatestGoodAnswer.Equal(stored.LatestGoodAnswer))
	assert.Equal(s.T(), projection.LatestGoodAnswerWasToday, stored.LatestGoodAnswerWasToday)
}
