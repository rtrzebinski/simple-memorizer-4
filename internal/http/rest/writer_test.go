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

func (suite *WriterSuite) TestStoreLesson() {
	lesson := models.Lesson{}

	method := "POST"
	route := StoreLesson
	params := map[string]string(nil)
	reqBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	suite.client.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.writer.StoreLesson(&lesson)
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

func (suite *WriterSuite) TestStoreExercise() {
	exercise := models.Exercise{}

	method := "POST"
	route := StoreExercise
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	suite.client.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.writer.StoreExercise(exercise)
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

func (suite *WriterSuite) TestIncrementBadAnswers() {
	exercise := models.Exercise{}

	method := "POST"
	route := IncrementBadAnswers
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	suite.client.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.writer.IncrementBadAnswers(exercise)
	suite.Assert().NoError(err)
}

func (suite *WriterSuite) TestIncrementGoodAnswers() {
	exercise := models.Exercise{}

	method := "POST"
	route := IncrementGoodAnswers
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	suite.client.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.writer.IncrementGoodAnswers(exercise)
	suite.Assert().NoError(err)
}
