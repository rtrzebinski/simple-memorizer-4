package api

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/routes"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/projections"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ClientSuite struct {
	suite.Suite
	s *Client
	c *CallerMock
}

func (suite *ClientSuite) SetupTest() {
	suite.c = NewCallerMock()
	suite.s = NewClient(suite.c)
}

func TestClientSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func (suite *ClientSuite) TestFetchLessons() {
	lessons := models.Lessons{models.Lesson{Name: "name"}}

	responseBody, err := json.Marshal(lessons)
	suite.Assert().NoError(err)

	method := "GET"
	route := routes.FetchLessons
	params := map[string]string(nil)
	reqBody := []byte(nil)

	suite.c.On("Call", method, route, params, reqBody).Return(responseBody)

	result, err := suite.s.FetchLessons()
	suite.Assert().NoError(err)
	suite.Assert().Equal(lessons, result)
}

func (suite *ClientSuite) TestHydrateLesson() {
	lesson := &models.Lesson{Id: 10}

	responseBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	method := "GET"
	route := routes.HydrateLesson
	params := map[string]string{"lesson_id": "10"}
	reqBody := []byte(nil)

	suite.c.On("Call", method, route, params, reqBody).Return(responseBody)

	err = suite.s.HydrateLesson(lesson)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestFetchExercises() {
	lesson := models.Lesson{Id: 1}
	exercises := models.Exercises{
		{Id: 1, Results: []models.Result{{Type: models.Good}, {Type: models.Good}}},
		{Id: 2, Results: []models.Result{{Type: models.Bad}, {Type: models.Bad}}},
	}

	expectedExercises := models.Exercises{
		{Id: 1, Results: []models.Result{{Type: models.Good}, {Type: models.Good}},
			ResultsProjection: projections.BuildResultsProjection(exercises[0].Results)},
		{Id: 2, Results: []models.Result{{Type: models.Bad}, {Type: models.Bad}},
			ResultsProjection: projections.BuildResultsProjection(exercises[1].Results)},
	}

	responseBody, err := json.Marshal(exercises)
	suite.Assert().NoError(err)

	method := "GET"
	route := routes.FetchExercises
	params := map[string]string{"lesson_id": "1"}
	reqBody := []byte(nil)

	suite.c.On("Call", method, route, params, reqBody).Return(responseBody)

	result, err := suite.s.FetchExercises(models.Lesson{Id: lesson.Id})

	suite.Nil(err)
	suite.Equal(expectedExercises, result)
	suite.Equal(expectedExercises[0].ResultsProjection, result[0].ResultsProjection)
	suite.Equal(expectedExercises[1].ResultsProjection, result[1].ResultsProjection)
	suite.c.AssertExpectations(suite.T())
}

func (suite *ClientSuite) TestUpsertLesson() {
	lesson := models.Lesson{}

	method := "POST"
	route := routes.UpsertLesson
	params := map[string]string(nil)
	reqBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	suite.c.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.s.UpsertLesson(&lesson)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestDeleteLesson() {
	lesson := models.Lesson{}

	method := "POST"
	route := routes.DeleteLesson
	params := map[string]string(nil)
	reqBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	suite.c.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.s.DeleteLesson(lesson)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestUpsertExercise() {
	exercise := models.Exercise{}

	method := "POST"
	route := routes.UpsertExercise
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	suite.c.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.s.UpsertExercise(&exercise)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestStoreExercises() {
	exercises := models.Exercises{}

	method := "POST"
	route := routes.StoreExercises
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercises)
	suite.Assert().NoError(err)

	suite.c.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.s.StoreExercises(exercises)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestDeleteExercise() {
	exercise := models.Exercise{}

	method := "POST"
	route := routes.DeleteExercise
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	suite.c.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.s.DeleteExercise(exercise)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestStoreResult() {
	result := models.Result{}

	method := "POST"
	route := routes.StoreResult
	params := map[string]string(nil)
	reqBody, err := json.Marshal(result)
	suite.Assert().NoError(err)

	suite.c.On("Call", method, route, params, reqBody).Return([]byte(""))

	err = suite.s.StoreResult(&result)
	suite.Assert().NoError(err)
}
