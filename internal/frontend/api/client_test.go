package api

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/server"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ClientSuite struct {
	suite.Suite
	client *Client
	caller *CallerMock
}

func (suite *ClientSuite) SetupTest() {
	suite.caller = NewCallerMock()
	suite.client = NewClient(suite.caller)
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func (suite *ClientSuite) TestClient_FetchLessons() {
	lessons := models.Lessons{models.Lesson{Name: "name"}}

	responseBody, err := json.Marshal(lessons)
	suite.Assert().NoError(err)

	method := "GET"
	route := server.FetchLessons
	params := map[string]string(nil)
	reqBody := []byte(nil)

	suite.caller.On("Call", method, route, params, reqBody).Return(responseBody)

	result, err := suite.client.FetchLessons()
	suite.Assert().NoError(err)
	suite.Assert().Equal(lessons, result)
}

func (suite *ClientSuite) TestClient_HydrateLesson() {
	lesson := &models.Lesson{Id: 10}

	responseBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	method := "GET"
	route := server.HydrateLesson
	params := map[string]string{"lesson_id": "10"}
	reqBody := []byte(nil)

	suite.caller.On("Call", method, route, params, reqBody).Return(responseBody)

	err = suite.client.HydrateLesson(lesson)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_FetchExercises() {
	lesson := models.Lesson{Id: 1}
	exercises := models.Exercises{
		{Id: 1, Results: []models.Result{{Type: models.Good}, {Type: models.Good}}},
		{Id: 2, Results: []models.Result{{Type: models.Bad}, {Type: models.Bad}}},
	}

	expectedExercises := models.Exercises{
		{Id: 1, Results: []models.Result{{Type: models.Good}, {Type: models.Good}},
			ResultsProjection: models.BuildResultsProjection(exercises[0].Results)},
		{Id: 2, Results: []models.Result{{Type: models.Bad}, {Type: models.Bad}},
			ResultsProjection: models.BuildResultsProjection(exercises[1].Results)},
	}

	responseBody, err := json.Marshal(exercises)
	suite.Assert().NoError(err)

	method := "GET"
	route := server.FetchExercises
	params := map[string]string{"lesson_id": "1"}
	reqBody := []byte(nil)

	suite.caller.On("Call", method, route, params, reqBody).Return(responseBody)

	result, err := suite.client.FetchExercises(models.Lesson{Id: lesson.Id})

	suite.Nil(err)
	suite.Equal(expectedExercises, result)
	suite.Equal(expectedExercises[0].ResultsProjection, result[0].ResultsProjection)
	suite.Equal(expectedExercises[1].ResultsProjection, result[1].ResultsProjection)
	suite.caller.AssertExpectations(suite.T())
}

func (suite *ClientSuite) TestClient_UpsertLesson() {
	lesson := models.Lesson{}

	method := "POST"
	route := server.UpsertLesson
	params := map[string]string(nil)
	reqBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	suite.caller.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.client.UpsertLesson(lesson)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_DeleteLesson() {
	lesson := models.Lesson{}

	method := "POST"
	route := server.DeleteLesson
	params := map[string]string(nil)
	reqBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	suite.caller.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.client.DeleteLesson(lesson)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_UpsertExercise() {
	exercise := models.Exercise{}

	method := "POST"
	route := server.UpsertExercise
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	suite.caller.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.client.UpsertExercise(exercise)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_StoreExercises() {
	exercises := models.Exercises{}

	method := "POST"
	route := server.StoreExercises
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercises)
	suite.Assert().NoError(err)

	suite.caller.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.client.StoreExercises(exercises)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_DeleteExercise() {
	exercise := models.Exercise{}

	method := "POST"
	route := server.DeleteExercise
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	suite.caller.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.client.DeleteExercise(exercise)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_StoreResult() {
	result := models.Result{}

	method := "POST"
	route := server.StoreResult
	params := map[string]string(nil)
	reqBody, err := json.Marshal(result)
	suite.Assert().NoError(err)

	suite.caller.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.client.StoreResult(result)
	suite.Assert().NoError(err)
}
