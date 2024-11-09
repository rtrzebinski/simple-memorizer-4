package worker

import (
	"context"
	"errors"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite
	writerMock *WriterMock
	service    *Service
}

func (suite *ServiceSuite) SetupTest() {
	suite.writerMock = new(WriterMock)
	suite.service = NewService(suite.writerMock)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (suite *ServiceSuite) TestService_ProcessGoodAnswer_Success() {
	exerciseID := 1
	result := &Result{
		Type:       Good,
		ExerciseId: exerciseID,
	}

	suite.writerMock.On("StoreResult", result).Return(nil)

	err := suite.service.ProcessGoodAnswer(context.Background(), exerciseID)

	suite.NoError(err)
	suite.writerMock.AssertCalled(suite.T(), "StoreResult", result)
}

func (suite *ServiceSuite) TestService_ProcessGoodAnswer_Error() {
	exerciseID := 1
	result := &Result{
		Type:       Good,
		ExerciseId: exerciseID,
	}

	suite.writerMock.On("StoreResult", result).Return(errors.New("database error"))

	err := suite.service.ProcessGoodAnswer(context.Background(), exerciseID)

	suite.Error(err)
	suite.Contains(err.Error(), "store good result")
	suite.writerMock.AssertCalled(suite.T(), "StoreResult", result)
}

func (suite *ServiceSuite) TestService_ProcessBadAnswer_Success() {
	exerciseID := 2
	result := &Result{
		Type:       Bad,
		ExerciseId: exerciseID,
	}

	suite.writerMock.On("StoreResult", result).Return(nil)

	err := suite.service.ProcessBadAnswer(context.Background(), exerciseID)

	suite.NoError(err)
	suite.writerMock.AssertCalled(suite.T(), "StoreResult", result)
}

func (suite *ServiceSuite) TestService_ProcessBadAnswer_Error() {
	exerciseID := 2
	result := &Result{
		Type:       Bad,
		ExerciseId: exerciseID,
	}

	suite.writerMock.On("StoreResult", result).Return(errors.New("database error"))

	err := suite.service.ProcessBadAnswer(context.Background(), exerciseID)

	suite.Error(err)
	suite.Contains(err.Error(), "store bad result")
	suite.writerMock.AssertCalled(suite.T(), "StoreResult", result)
}
