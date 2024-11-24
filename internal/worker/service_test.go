package worker

import (
	"context"
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite
	r       *ReaderMock
	w       *WriterMock
	pb      *ProjectionBuilderMock
	service *Service
}

func (suite *ServiceSuite) SetupTest() {
	suite.r = new(ReaderMock)
	suite.w = new(WriterMock)
	suite.pb = new(ProjectionBuilderMock)
	suite.service = NewService(suite.r, suite.w, suite.pb)
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

	projection := ResultsProjection{
		GoodAnswers: 1,
	}

	suite.w.On("StoreResult", ctx, result).Return(nil)
	suite.r.On("FetchResults", ctx, exerciseID).Return(results, nil)
	suite.pb.On("Projection", results).Return(projection)
	suite.w.On("UpdateExerciseProjection", ctx, exerciseID, projection).Return(nil)

	err := suite.service.ProcessGoodAnswer(ctx, exerciseID)

	suite.NoError(err)
	suite.r.AssertExpectations(suite.T())
	suite.w.AssertExpectations(suite.T())
	suite.pb.AssertExpectations(suite.T())
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
	suite.pb.AssertExpectations(suite.T())
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

	projection := ResultsProjection{
		BadAnswers: 1,
	}

	suite.w.On("StoreResult", ctx, result).Return(nil)
	suite.r.On("FetchResults", ctx, exerciseID).Return(results, nil)
	suite.pb.On("Projection", results).Return(projection)
	suite.w.On("UpdateExerciseProjection", ctx, exerciseID, projection).Return(nil)

	err := suite.service.ProcessBadAnswer(ctx, exerciseID)

	suite.NoError(err)
	suite.r.AssertExpectations(suite.T())
	suite.w.AssertExpectations(suite.T())
	suite.pb.AssertExpectations(suite.T())
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
	suite.pb.AssertExpectations(suite.T())
}
