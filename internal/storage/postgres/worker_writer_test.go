package postgres

import (
	"testing"
	"time"

	"github.com/guregu/null/v5"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/worker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type WorkerWriterSuite struct {
	Suite
	writer *WorkerWriter
}

func TestWorkerWriter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	suite.Run(t, new(WorkerWriterSuite))
}

func (s *WorkerWriterSuite) SetupSuite() {
	s.Suite.SetupSuite()
	s.writer = NewWorkerWriter(s.DB)
}
func (s *WorkerWriterSuite) TestWorkerWriter_StoreAnswer() {
	ctx := s.T().Context()

	exercise := &exercise{}
	createExercise(s.DB, exercise)

	result := worker.Result{
		Type:       worker.Good,
		ExerciseId: exercise.Id,
	}

	err := s.writer.StoreResult(ctx, result)
	assert.NoError(s.T(), err)

	stored := fetchLatestResult(s.DB)

	assert.Equal(s.T(), string(result.Type), stored.Type)
	assert.Equal(s.T(), result.ExerciseId, stored.ExerciseId)
}

func (s *WorkerWriterSuite) TestWorkerWriter_UpdateExerciseProjection_allProjections() {
	ctx := s.T().Context()

	exercise := &exercise{}
	createExercise(s.DB, exercise)

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

	err := s.writer.UpdateExerciseProjection(ctx, exercise.Id, projection)
	assert.NoError(s.T(), err)

	stored := findExerciseById(s.DB, exercise.Id)

	assert.Equal(s.T(), projection.BadAnswers, stored.BadAnswers)
	assert.Equal(s.T(), projection.BadAnswersToday, stored.BadAnswersToday)
	assert.True(s.T(), projection.LatestBadAnswer.Equal(stored.LatestBadAnswer))
	assert.Equal(s.T(), projection.LatestBadAnswerWasToday, stored.LatestBadAnswerWasToday)
	assert.Equal(s.T(), projection.GoodAnswers, stored.GoodAnswers)
	assert.Equal(s.T(), projection.GoodAnswersToday, stored.GoodAnswersToday)
	assert.True(s.T(), projection.LatestGoodAnswer.Equal(stored.LatestGoodAnswer))
	assert.Equal(s.T(), projection.LatestGoodAnswerWasToday, stored.LatestGoodAnswerWasToday)
}

func (s *WorkerWriterSuite) TestWorkerWriter_UpdateExerciseProjection_badOnly() {
	ctx := s.T().Context()

	exercise := &exercise{}
	createExercise(s.DB, exercise)

	projection := worker.ResultsProjection{
		BadAnswers:              1,
		BadAnswersToday:         2,
		LatestBadAnswer:         null.TimeFrom(time.Date(2021, time.March, 1, 12, 30, 0, 0, time.UTC)),
		LatestBadAnswerWasToday: true,
	}

	err := s.writer.UpdateExerciseProjection(ctx, exercise.Id, projection)
	assert.NoError(s.T(), err)

	stored := findExerciseById(s.DB, exercise.Id)

	assert.Equal(s.T(), projection.BadAnswers, stored.BadAnswers)
	assert.Equal(s.T(), projection.BadAnswersToday, stored.BadAnswersToday)
	assert.True(s.T(), projection.LatestBadAnswer.Equal(stored.LatestBadAnswer))
	assert.Equal(s.T(), projection.LatestBadAnswerWasToday, stored.LatestBadAnswerWasToday)
	assert.Equal(s.T(), projection.GoodAnswers, stored.GoodAnswers)
	assert.Equal(s.T(), projection.GoodAnswersToday, stored.GoodAnswersToday)
	assert.True(s.T(), projection.LatestGoodAnswer.Equal(stored.LatestGoodAnswer))
	assert.Equal(s.T(), projection.LatestGoodAnswerWasToday, stored.LatestGoodAnswerWasToday)
}

func (s *WorkerWriterSuite) TestWorkerWriter_UpdateExerciseProjection_goodOnly() {
	ctx := s.T().Context()

	exercise := &exercise{}
	createExercise(s.DB, exercise)

	projection := worker.ResultsProjection{
		GoodAnswers:              3,
		GoodAnswersToday:         4,
		LatestGoodAnswer:         null.TimeFrom(time.Date(2022, time.March, 1, 12, 30, 0, 0, time.UTC)),
		LatestGoodAnswerWasToday: true,
	}

	err := s.writer.UpdateExerciseProjection(ctx, exercise.Id, projection)
	assert.NoError(s.T(), err)

	stored := findExerciseById(s.DB, exercise.Id)

	assert.Equal(s.T(), projection.BadAnswers, stored.BadAnswers)
	assert.Equal(s.T(), projection.BadAnswersToday, stored.BadAnswersToday)
	assert.True(s.T(), projection.LatestBadAnswer.Equal(stored.LatestBadAnswer))
	assert.Equal(s.T(), projection.LatestBadAnswerWasToday, stored.LatestBadAnswerWasToday)
	assert.Equal(s.T(), projection.GoodAnswers, stored.GoodAnswers)
	assert.Equal(s.T(), projection.GoodAnswersToday, stored.GoodAnswersToday)
	assert.True(s.T(), projection.LatestGoodAnswer.Equal(stored.LatestGoodAnswer))
	assert.Equal(s.T(), projection.LatestGoodAnswerWasToday, stored.LatestGoodAnswerWasToday)
}
