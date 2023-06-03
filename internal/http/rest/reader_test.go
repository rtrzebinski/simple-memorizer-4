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

type ReaderSuite struct {
	suite.Suite

	http   *myhttp.DoerMock
	reader *Reader
}

func (suite *ReaderSuite) SetupTest() {
	suite.http = new(myhttp.DoerMock)
	suite.reader = NewReader(suite.http, host, scheme)
}

func TestReaderSuite(t *testing.T) {
	suite.Run(t, new(ReaderSuite))
}

func (suite *ReaderSuite) TestFetchAllLessons() {
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

	result, err := suite.reader.FetchAllLessons()
	assert.NoError(suite.T(), err)

	suite.Assert().Equal(lessons, result)
}

func (suite *ReaderSuite) TestFetchExercisesOfLesson() {
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

	result, err := suite.reader.FetchExercisesOfLesson(models.Lesson{Id: lessonId})
	assert.NoError(suite.T(), err)

	suite.Assert().Equal(exercises, result)
}

func (suite *ReaderSuite) TestRandomExerciseOfLesson() {
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
		suite.Equal(FetchRandomExerciseOfLesson+"?lesson_id=10", req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		lId, _ := strconv.Atoi(req.URL.Query().Get("lesson_id"))
		suite.Equal(lessonId, lId)

		return true

	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	result, err := suite.reader.FetchRandomExerciseOfLesson(models.Lesson{Id: lessonId})
	assert.NoError(suite.T(), err)

	suite.Assert().Equal(exercise, result)
}
