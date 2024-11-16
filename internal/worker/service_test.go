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

	suite.w.On("StoreResult", result).Return(nil)
	suite.r.On("FetchResults", exerciseID).Return(results, nil)
	suite.pb.On("Projection", results).Return(projection)
	suite.w.On("UpdateExerciseProjection", exerciseID, projection).Return(nil)

	err := suite.service.ProcessGoodAnswer(context.Background(), exerciseID)

	suite.NoError(err)
	suite.r.AssertExpectations(suite.T())
	suite.w.AssertExpectations(suite.T())
	suite.pb.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestService_ProcessGoodAnswer_Error() {
	exerciseID := 1
	result := Result{
		Type:       Good,
		ExerciseId: exerciseID,
	}

	suite.w.On("StoreResult", result).Return(errors.New("database error"))

	err := suite.service.ProcessGoodAnswer(context.Background(), exerciseID)

	suite.Error(err)
	suite.Contains(err.Error(), "database error")
	suite.r.AssertExpectations(suite.T())
	suite.w.AssertExpectations(suite.T())
	suite.pb.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestService_ProcessBadAnswer_Success() {
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

	suite.w.On("StoreResult", result).Return(nil)
	suite.r.On("FetchResults", exerciseID).Return(results, nil)
	suite.pb.On("Projection", results).Return(projection)
	suite.w.On("UpdateExerciseProjection", exerciseID, projection).Return(nil)

	err := suite.service.ProcessBadAnswer(context.Background(), exerciseID)

	suite.NoError(err)
	suite.r.AssertExpectations(suite.T())
	suite.w.AssertExpectations(suite.T())
	suite.pb.AssertExpectations(suite.T())
}

func (suite *ServiceSuite) TestService_ProcessBadAnswer_Error() {
	exerciseID := 2
	result := Result{
		Type:       Bad,
		ExerciseId: exerciseID,
	}

	suite.w.On("StoreResult", result).Return(errors.New("database error"))

	err := suite.service.ProcessBadAnswer(context.Background(), exerciseID)

	suite.Error(err)
	suite.Contains(err.Error(), "database error")
	suite.r.AssertExpectations(suite.T())
	suite.w.AssertExpectations(suite.T())
	suite.pb.AssertExpectations(suite.T())
}
