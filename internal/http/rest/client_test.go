package rest

import (
	"bytes"
	"encoding/json"
	myhttp "github.com/rtrzebinski/simple-memorizer-4/internal/http"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"

	"io"
	"net/http"
	"strconv"
	"testing"
)

var (
	host   = "example.com"
	scheme = "http"
)

type ClientSuite struct {
	suite.Suite

	http   *myhttp.DoerMock
	client *Client
}

func (suite *ClientSuite) SetupTest() {
	suite.http = new(myhttp.DoerMock)
	suite.client = NewClient(suite.http, host, scheme)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ClientSuite))
}

func (suite *ClientSuite) TestDeleteExercise() {
	exercise := models.Exercise{
		Id: 123,
	}

	suite.http.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("POST", req.Method)
		suite.Equal(DeleteExercise, req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		var input models.Exercise
		err := json.NewDecoder(req.Body).Decode(&input)
		suite.NoError(err)
		suite.Equal(exercise.Question, input.Question)
		suite.Equal(exercise.Answer, input.Answer)

		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}, nil)

	err := suite.client.DeleteExercise(exercise)
	assert.NoError(suite.T(), err)
}

func (suite *ClientSuite) TestStoreExercise() {
	exercise := models.Exercise{
		Question: "question",
		Answer:   "answer",
	}

	suite.http.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("POST", req.Method)
		suite.Equal(StoreExercise, req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		var input models.Exercise
		err := json.NewDecoder(req.Body).Decode(&input)
		suite.NoError(err)
		suite.Equal(exercise.Question, input.Question)
		suite.Equal(exercise.Answer, input.Answer)

		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}, nil)

	err := suite.client.StoreExercise(exercise)
	assert.NoError(suite.T(), err)
}

func (suite *ClientSuite) TestFetchExercisesOfLesson() {
	exercise := models.Exercise{
		Id:          1,
		Question:    "question",
		Answer:      "answer",
		BadAnswers:  2,
		GoodAnswers: 3,
	}
	exercises := models.Exercises{exercise}

	lessonId := 10

	responseBody, err := json.Marshal(exercises)
	if err != nil {
		suite.Error(err)
	}

	suite.http.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("GET", req.Method)
		suite.Equal(FetchExercisesOfLesson+"?lesson_id=10", req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		lId, _ := strconv.Atoi(req.URL.Query().Get("lesson_id"))
		suite.Equal(lessonId, lId)

		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	result, err := suite.client.FetchExercisesOfLesson(lessonId)
	assert.NoError(suite.T(), err)

	suite.Assert().Equal(exercises, result)
}

func (suite *ClientSuite) TestDeleteLeson() {
	lesson := models.Lesson{
		Id: 123,
	}

	suite.http.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("POST", req.Method)
		suite.Equal(DeleteLesson, req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		var input models.Lesson
		err := json.NewDecoder(req.Body).Decode(&input)
		suite.NoError(err)
		suite.Equal(lesson.Name, input.Name)

		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}, nil)

	err := suite.client.DeleteLesson(lesson)
	assert.NoError(suite.T(), err)
}

func (suite *ClientSuite) TestStoreLesson() {
	lesson := models.Lesson{
		Name: "name",
	}

	suite.http.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("POST", req.Method)
		suite.Equal(StoreLesson, req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		var input models.Lesson
		err := json.NewDecoder(req.Body).Decode(&input)
		suite.NoError(err)
		suite.Equal(lesson.Name, input.Name)

		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}, nil)

	err := suite.client.StoreLesson(lesson)
	assert.NoError(suite.T(), err)
}

func (suite *ClientSuite) TestFetchAllLessons() {
	lesson := models.Lesson{
		Id:   1,
		Name: "name",
	}
	lessons := models.Lessons{lesson}

	responseBody, err := json.Marshal(lessons)
	if err != nil {
		suite.Error(err)
	}

	suite.http.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("GET", req.Method)
		suite.Equal(FetchAllLessons, req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	result, err := suite.client.FetchAllLessons()
	assert.NoError(suite.T(), err)

	suite.Assert().Equal(lessons, result)
}

func (suite *ClientSuite) TestFetchNextExerciseOfLesson() {
	exercise := models.Exercise{
		Id:          1,
		Question:    "question",
		Answer:      "answer",
		BadAnswers:  2,
		GoodAnswers: 3,
	}

	responseBody, err := json.Marshal(exercise)
	if err != nil {
		suite.Error(err)
	}

	lessonId := 10

	suite.http.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("GET", req.Method)
		suite.Equal(FetchNextExerciseOfLesson+"?lesson_id=10", req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		lId, _ := strconv.Atoi(req.URL.Query().Get("lesson_id"))
		suite.Equal(lessonId, lId)

		return true

	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	result, err := suite.client.FetchNextExerciseOfLesson(lessonId)
	assert.NoError(suite.T(), err)

	suite.Assert().Equal(exercise, result)
}

func (suite *ClientSuite) TestIncrementBadAnswers() {
	exercise := models.Exercise{Id: 123}

	suite.http.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("POST", req.Method)
		suite.Equal(IncrementBadAnswers, req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		var input models.Exercise
		err := json.NewDecoder(req.Body).Decode(&input)
		suite.NoError(err)
		suite.Equal(exercise.Id, input.Id)

		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}, nil)

	err := suite.client.IncrementBadAnswers(exercise)
	assert.NoError(suite.T(), err)
}

func (suite *ClientSuite) TestIncrementGoodAnswers() {
	exercise := models.Exercise{Id: 123}

	suite.http.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("POST", req.Method)
		suite.Equal(IncrementGoodAnswers, req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		var input models.Exercise
		err := json.NewDecoder(req.Body).Decode(&input)
		suite.NoError(err)
		suite.Equal(exercise.Id, input.Id)

		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}, nil)

	err := suite.client.IncrementGoodAnswers(exercise)
	assert.NoError(suite.T(), err)
}
