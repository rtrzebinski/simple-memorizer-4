package frontend

import (
	"bytes"
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/routes"
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

func (suite *ApiClientSuite) TestFetchRandomExercise() {
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
		suite.Equal(backend.FetchRandomExercise, req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	result, err := suite.apiClient.FetchRandomExercise()
	assert.NoError(suite.T(), err)

	suite.Assert().Equal(exercise, result)
}

func (suite *ApiClientSuite) TestIncrementBadAnswers() {
	exerciseId := 123

	suite.httpClientMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("POST", req.Method)
		suite.Equal(backend.IncrementBadAnswers, req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		var input routes.IncrementBadAnswersReq
		err := json.NewDecoder(req.Body).Decode(&input)
		suite.NoError(err)
		suite.Equal(exerciseId, input.ExerciseId)

		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}, nil)

	err := suite.apiClient.IncrementBadAnswers(exerciseId)
	assert.NoError(suite.T(), err)
}

func (suite *ApiClientSuite) TestIncrementGoodAnswers() {
	exerciseId := 456

	suite.httpClientMock.On("Do", mock.MatchedBy(func(req *http.Request) bool {
		suite.Equal("POST", req.Method)
		suite.Equal(backend.IncrementGoodAnswers, req.URL.RequestURI())
		suite.Equal(host, req.URL.Host)
		suite.Equal(scheme, req.URL.Scheme)

		var input routes.IncrementGoodAnswersReq
		err := json.NewDecoder(req.Body).Decode(&input)
		suite.NoError(err)
		suite.Equal(exerciseId, input.ExerciseId)

		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader([]byte{})),
	}, nil)

	err := suite.apiClient.IncrementGoodAnswers(exerciseId)
	assert.NoError(suite.T(), err)
}
