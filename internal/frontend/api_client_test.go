package frontend

import (
	"bytes"
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"io"
	"net/http"
	"testing"
)

type ApiClientSuite struct {
	suite.Suite

	httpClientMock *HttpClientMock
	apiClient      *ApiClient
}

func (suite *ApiClientSuite) SetupTest() {
	suite.httpClientMock = new(HttpClientMock)
	suite.apiClient = NewApiClient(suite.httpClientMock, "example.com", "http")
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
		return true
	})).Return(&http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(responseBody)),
	}, nil)

	result := suite.apiClient.FetchRandomExercise()

	suite.Assert().Equal(exercise, result)
}

func (suite *ApiClientSuite) TestIncrementBadAnswers() {

}

func (suite *ApiClientSuite) TestIncrementGoodAnswers() {

}
