package api

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/routes"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/projections"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ServiceSuite struct {
	suite.Suite
	s *Service
	c *CallerMock
}

func (suite *ServiceSuite) SetupTest() {
	suite.c = NewCallerMock()
	suite.s = NewService(suite.c)
}

func TestServiceSuite(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}

func (suite *ServiceSuite) TestFetchLessons() {
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

func (suite *ServiceSuite) TestHydrateLesson() {
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

func (suite *ServiceSuite) TestFetchExercises() {
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

func (suite *ServiceSuite) TestUpsertLesson() {
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

func (suite *ServiceSuite) TestDeleteLesson() {
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

func (suite *ServiceSuite) TestUpsertExercise() {
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

func (suite *ServiceSuite) TestStoreExercises() {
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

func (suite *ServiceSuite) TestDeleteExercise() {
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

func (suite *ServiceSuite) TestStoreResult() {
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
