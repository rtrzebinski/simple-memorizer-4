package rest

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/routes"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ReaderSuite struct {
	suite.Suite

	client *frontend.ClientMock
	reader *Reader
}

func (suite *ReaderSuite) SetupTest() {
	suite.client = frontend.NewClientMock()
	suite.reader = NewReader(suite.client)
}

func TestReaderSuite(t *testing.T) {
	suite.Run(t, new(ReaderSuite))
}

func (suite *ReaderSuite) TestFetchLessons() {
	lessons := models.Lessons{models.Lesson{Name: "name"}}

	responseBody, err := json.Marshal(lessons)
	suite.Assert().NoError(err)

	method := "GET"
	route := routes.FetchLessons
	params := map[string]string(nil)
	reqBody := []byte(nil)

	suite.client.On("Call", method, route, params, reqBody).Return(responseBody)

	result, err := suite.reader.FetchLessons()
	suite.Assert().NoError(err)
	suite.Assert().Equal(lessons, result)
}

func (suite *ReaderSuite) TestHydrateLesson() {
	lesson := &models.Lesson{Id: 10}

	responseBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	method := "GET"
	route := routes.HydrateLesson
	params := map[string]string{"lesson_id": "10"}
	reqBody := []byte(nil)

	suite.client.On("Call", method, route, params, reqBody).Return(responseBody)

	err = suite.reader.HydrateLesson(lesson)
	suite.Assert().NoError(err)
}

func (suite *ReaderSuite) TestFetchExercises() {
	exercises := models.Exercises{models.Exercise{Question: "question"}}
	lessonId := 10

	responseBody, err := json.Marshal(exercises)
	suite.Assert().NoError(err)

	method := "GET"
	route := routes.FetchExercises
	params := map[string]string{"lesson_id": "10"}
	reqBody := []byte(nil)

	suite.client.On("Call", method, route, params, reqBody).Return(responseBody)

	result, err := suite.reader.FetchExercises(models.Lesson{Id: lessonId})

	suite.Assert().NoError(err)
	suite.Assert().Equal(exercises, result)
}
