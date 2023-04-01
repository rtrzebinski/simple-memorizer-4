package frontend

import (
	"bytes"
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"testing"
)

var (
	host   = "example.com"
	scheme = "http"
)

type ApiClientSuite struct {
	suite.Suite

	httpClientMock *HttpClientMock
	apiClient      *ApiClient
}

func (suite *ApiClientSuite) SetupTest() {
	suite.httpClientMock = new(HttpClientMock)
	suite.apiClient = NewApiClient(suite.httpClientMock, host, scheme)
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(ApiClientSuite))
}

func (suite *ApiClientSuite) TestStoreExercise() {
	exercise := models.Exercise{
		Question: "question",
		Answer:   "answer",
	}

	suite.httpClientMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("POST", req.Method)
		suite.Equal(backend.StoreExercise, req.URL.RequestURI())
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

	err := suite.apiClient.StoreExercise(exercise)
	assert.NoError(suite.T(), err)
}

func (suite *ApiClientSuite) TestFetchExercises() {
	exercise := models.Exercise{
		Id:          1,
		Question:    "question",
		Answer:      "answer",
		BadAnswers:  2,
		GoodAnswers: 3,
	}
	exercises := models.Exercises{exercise}

	responseBody, err := json.Marshal(exercises)
	if err != nil {
		suite.Error(err)
	}

	suite.httpClientMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("GET", req.Method)
		suite.Equal(backend.FetchExercises, req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	result, err := suite.apiClient.FetchExercises()
	assert.NoError(suite.T(), err)

	suite.Assert().Equal(exercises, result)
}

func (suite *ApiClientSuite) TestFetchNextExercise() {
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

	suite.httpClientMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("GET", req.Method)
		suite.Equal(backend.FetchNextExercise, req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	result, err := suite.apiClient.FetchNextExercise()
	assert.NoError(suite.T(), err)

	suite.Assert().Equal(exercise, result)
}

func (suite *ApiClientSuite) TestIncrementBadAnswers() {
	exercise := models.Exercise{Id: 123}

	suite.httpClientMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("POST", req.Method)
		suite.Equal(backend.IncrementBadAnswers, req.URL.RequestURI())
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

	err := suite.apiClient.IncrementBadAnswers(exercise)
	assert.NoError(suite.T(), err)
}

func (suite *ApiClientSuite) TestIncrementGoodAnswers() {
	exercise := models.Exercise{Id: 123}

	suite.httpClientMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("POST", req.Method)
		suite.Equal(backend.IncrementGoodAnswers, req.URL.RequestURI())
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

	err := suite.apiClient.IncrementGoodAnswers(exercise)
	assert.NoError(suite.T(), err)
}
