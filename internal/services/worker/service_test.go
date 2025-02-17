package worker

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type ServiceSuite struct {
	suite.Suite
	r       *ReaderMock
	w       *WriterMock
	service *Service
}

func (suite *ServiceSuite) SetupTest() {
	suite.r = new(ReaderMock)
	suite.w = new(WriterMock)
	suite.service = NewService(suite.r, suite.w)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (suite *ServiceSuite) TestService_ProcessGoodAnswer_Success() {
	ctx := context.Background()

	exerciseID := 1
	result := Result{
		Type:       Good,
		ExerciseId: exerciseID,
	}

	results := []Result{
		{
			Type:       Good,
			ExerciseId: exerciseID,
		},
	}

	rp := ResultsProjection{
		GoodAnswers: 1,
	}

	suite.w.On("StoreResult", ctx, result).Return(nil)
	suite.r.On("FetchResults", ctx, exerciseID).Return(results, nil)
	suite.w.On("UpdateExerciseProjection", ctx, exerciseID, rp).Return(nil)

	err := suite.service.ProcessGoodAnswer(ctx, exerciseID)

	suite.NoError(err)
	suite.r.AssertExpectations(suite.T())
	suite.w.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestService_ProcessGoodAnswer_Error() {
	ctx := context.Background()

	exerciseID := 1
	result := Result{
		Type:       Good,
		ExerciseId: exerciseID,
	}

	suite.w.On("StoreResult", ctx, result).Return(errors.New("database error"))

	err := suite.service.ProcessGoodAnswer(ctx, exerciseID)

	suite.Error(err)
	suite.Contains(err.Error(), "database error")
	suite.r.AssertExpectations(suite.T())
	suite.w.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestService_ProcessBadAnswer_Success() {
	ctx := context.Background()

	exerciseID := 2
	result := Result{
		Type:       Bad,
		ExerciseId: exerciseID,
	}

	results := []Result{
		{
			Type:       Bad,
			ExerciseId: exerciseID,
		},
	}

	rp := ResultsProjection{
		BadAnswers: 1,
	}

	suite.w.On("StoreResult", ctx, result).Return(nil)
	suite.r.On("FetchResults", ctx, exerciseID).Return(results, nil)
	suite.w.On("UpdateExerciseProjection", ctx, exerciseID, rp).Return(nil)

	err := suite.service.ProcessBadAnswer(ctx, exerciseID)

	suite.NoError(err)
	suite.r.AssertExpectations(suite.T())
	suite.w.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestService_ProcessBadAnswer_Error() {
	ctx := context.Background()

	exerciseID := 2
	result := Result{
		Type:       Bad,
		ExerciseId: exerciseID,
	}

	suite.w.On("StoreResult", ctx, result).Return(errors.New("database error"))

	err := suite.service.ProcessBadAnswer(ctx, exerciseID)

	suite.Error(err)
	suite.Contains(err.Error(), "database error")
	suite.r.AssertExpectations(suite.T())
	suite.w.AssertExpectations(suite.T())
}

func TestService_resultsProjection(t *testing.T) {
	yesterday := time.Now().Add(-24 * time.Hour)
	today := time.Now()

	var results []Result

	results = append(results, Result{
		Type:      Bad,
		CreatedAt: yesterday,
	})

	results = append(results, Result{
		Type:      Bad,
		CreatedAt: yesterday,
	})

	results = append(results, Result{
		Type:      Bad,
		CreatedAt: yesterday,
	})

	results = append(results, Result{
		Type:      Good,
		CreatedAt: yesterday,
	})

	rp := resultsProjection(results)

	assert.Equal(t, 3, rp.BadAnswers)
	assert.Equal(t, 0, rp.BadAnswersToday)
	assert.Equal(t, yesterday, rp.LatestBadAnswer.Time)
	assert.False(t, rp.LatestBadAnswerWasToday)
	assert.Equal(t, 1, rp.GoodAnswers)
	assert.Equal(t, 0, rp.GoodAnswersToday)
	assert.Equal(t, yesterday, rp.LatestGoodAnswer.Time)
	assert.False(t, rp.LatestGoodAnswerWasToday)

	results = append(results, Result{
		Type:      Bad,
		CreatedAt: today,
	})

	results = append(results, Result{
		Type:      Bad,
		CreatedAt: today,
	})

	results = append(results, Result{
		Type:      Good,
		CreatedAt: today,
	})

	results = append(results, Result{
		Type:      Good,
		CreatedAt: today,
	})

	rp = resultsProjection(results)

	assert.Equal(t, 5, rp.BadAnswers)
	assert.Equal(t, 2, rp.BadAnswersToday)
	assert.Equal(t, today, rp.LatestBadAnswer.Time)
	assert.True(t, rp.LatestBadAnswerWasToday)
	assert.Equal(t, 3, rp.GoodAnswers)
	assert.Equal(t, 2, rp.GoodAnswersToday)
	assert.Equal(t, today, rp.LatestGoodAnswer.Time)
	assert.True(t, rp.LatestGoodAnswerWasToday)
}
