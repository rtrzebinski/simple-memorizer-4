package rest

import (
	"encoding/json"
	myhttp "github.com/rtrzebinski/simple-memorizer-4/internal/http"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ReaderSuite struct {
	suite.Suite

	client *myhttp.ClientMock
	reader *Reader
}

func (suite *ReaderSuite) SetupTest() {
	suite.client = myhttp.NewClientMock()
	suite.reader = NewReader(suite.client)
}

func TestReaderSuite(t *testing.T) {
	suite.Run(t, new(ReaderSuite))
}

func (suite *ReaderSuite) TestFetchAllLessons() {
	lessons := models.Lessons{}

	responseBody, err := json.Marshal(lessons)
	suite.Assert().NoError(err)

	method := "GET"
	route := FetchAllLessons
	params := map[string]string(nil)
	reqBody := []byte(nil)

	suite.client.On("Call", method, route, params, reqBody).Return(responseBody)

	result, err := suite.reader.FetchAllLessons()
	suite.Assert().NoError(err)
	suite.Assert().Equal(lessons, result)
}

func (suite *ReaderSuite) TestHydrateLesson() {
	lesson := &models.Lesson{Id: 10}

	responseBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	method := "GET"
	route := HydrateLesson
	params := map[string]string{"lesson_id": "10"}
	reqBody := []byte(nil)

	suite.client.On("Call", method, route, params, reqBody).Return(responseBody)

	err = suite.reader.HydrateLesson(lesson)
	suite.Assert().NoError(err)
}

func (suite *ReaderSuite) TestFetchExercisesOfLesson() {
	exercises := models.Exercises{}
	lessonId := 10

	responseBody, err := json.Marshal(exercises)
	suite.Assert().NoError(err)

	method := "GET"
	route := FetchExercisesOfLesson
	params := map[string]string{"lesson_id": "10"}
	reqBody := []byte(nil)

	suite.client.On("Call", method, route, params, reqBody).Return(responseBody)

	result, err := suite.reader.FetchExercisesOfLesson(models.Lesson{Id: lessonId})

	suite.Assert().NoError(err)
	suite.Assert().Equal(exercises, result)
}

func (suite *ReaderSuite) TestRandomExerciseOfLesson() {
	exercise := models.Exercise{}
	lessonId := 10

	responseBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	method := "GET"
	route := FetchRandomExerciseOfLesson
	params := map[string]string{"lesson_id": "10"}
	reqBody := []byte(nil)

	suite.client.On("Call", method, route, params, reqBody).Return(responseBody)

	result, err := suite.reader.FetchRandomExerciseOfLesson(models.Lesson{Id: lessonId})

	suite.Assert().NoError(err)
	suite.Assert().Equal(exercise, result)
}
