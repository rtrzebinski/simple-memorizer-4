package rest

import (
	"encoding/json"
	myhttp "github.com/rtrzebinski/simple-memorizer-4/internal/http"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/suite"
	"testing"
)

type WriterSuite struct {
	suite.Suite

	client *myhttp.ClientMock
	writer *Writer
}

func (suite *WriterSuite) SetupTest() {
	suite.client = myhttp.NewClientMock()
	suite.writer = NewWriter(suite.client)
}

func TestWriterSuite(t *testing.T) {
	suite.Run(t, new(WriterSuite))
}

func (suite *WriterSuite) TestUpsertLesson() {
	lesson := models.Lesson{}

	method := "POST"
	route := UpsertLesson
	params := map[string]string(nil)
	reqBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	suite.client.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.writer.UpsertLesson(&lesson)
	suite.Assert().NoError(err)
}

func (suite *WriterSuite) TestDeleteLesson() {
	lesson := models.Lesson{}

	method := "POST"
	route := DeleteLesson
	params := map[string]string(nil)
	reqBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	suite.client.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.writer.DeleteLesson(lesson)
	suite.Assert().NoError(err)
}

func (suite *WriterSuite) TestUpsertExercise() {
	exercise := models.Exercise{}

	method := "POST"
	route := UpsertExercise
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	suite.client.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.writer.UpsertExercise(&exercise)
	suite.Assert().NoError(err)
}

func (suite *WriterSuite) TestStoreExercises() {
	exercises := models.Exercises{}

	method := "POST"
	route := StoreExercises
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercises)
	suite.Assert().NoError(err)

	suite.client.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.writer.StoreExercises(exercises)
	suite.Assert().NoError(err)
}

func (suite *WriterSuite) TestDeleteExercise() {
	exercise := models.Exercise{}

	method := "POST"
	route := DeleteExercise
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	suite.client.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.writer.DeleteExercise(exercise)
	suite.Assert().NoError(err)
}

func (suite *WriterSuite) TestStoreResult() {
	result := models.Result{}

	method := "POST"
	route := StoreResult
	params := map[string]string(nil)
	reqBody, err := json.Marshal(result)
	suite.Assert().NoError(err)

	suite.client.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.writer.StoreResult(&result)
	suite.Assert().NoError(err)
}
