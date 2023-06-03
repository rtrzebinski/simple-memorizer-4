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
	"testing"
)

// todo move to suite
//var (
//	host   = "example.com"
//	scheme = "http"
//)

type WriterSuite struct {
	suite.Suite

	http   *myhttp.DoerMock
	writer *Writer
}

func (suite *WriterSuite) SetupTest() {
	suite.http = new(myhttp.DoerMock)
	suite.writer = NewWriter(suite.http, host, scheme)
}

func TestWriterSuite(t *testing.T) {
	suite.Run(t, new(WriterSuite))
}

func (suite *WriterSuite) TestStoreLesson() {
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

	err := suite.writer.StoreLesson(lesson)
	assert.NoError(suite.T(), err)
}

func (suite *WriterSuite) TestDeleteLesson() {
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

	err := suite.writer.DeleteLesson(lesson)
	assert.NoError(suite.T(), err)
}

func (suite *WriterSuite) TestStoreExercise() {
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

	err := suite.writer.StoreExercise(exercise)
	assert.NoError(suite.T(), err)
}

func (suite *WriterSuite) TestDeleteExercise() {
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

	err := suite.writer.DeleteExercise(exercise)
	assert.NoError(suite.T(), err)
}

func (suite *WriterSuite) TestIncrementBadAnswers() {
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

	err := suite.writer.IncrementBadAnswers(exercise)
	assert.NoError(suite.T(), err)
}

func (suite *WriterSuite) TestIncrementGoodAnswers() {
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

	err := suite.writer.IncrementGoodAnswers(exercise)
	assert.NoError(suite.T(), err)
}
