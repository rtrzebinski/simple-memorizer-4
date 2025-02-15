package api

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/guregu/null/v5"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/server"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
	"github.com/stretchr/testify/suite"
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
	ctx := context.Background()

	lessons := []frontend.Lesson{{Name: "name"}}

	responseBody, err := json.Marshal(lessons)
	suite.Assert().NoError(err)

	method := "GET"
	route := server.FetchLessons
	params := map[string]string(nil)
	reqBody := []byte(nil)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return(responseBody)

	result, err := suite.client.FetchLessons(ctx)
	suite.Assert().NoError(err)
	suite.Assert().Equal(lessons, result)
}

func (suite *ClientSuite) TestClient_HydrateLesson() {
	ctx := context.Background()

	lesson := &frontend.Lesson{Id: 10}

	responseBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	method := "GET"
	route := server.HydrateLesson
	params := map[string]string{"lesson_id": "10"}
	reqBody := []byte(nil)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return(responseBody)

	err = suite.client.HydrateLesson(ctx, lesson)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_FetchExercises() {
	ctx := context.Background()

	lesson := frontend.Lesson{Id: 1}
	exercises := []frontend.Exercise{
		{
			Id:                       1,
			Question:                 "question",
			Answer:                   "answer",
			BadAnswers:               2,
			BadAnswersToday:          1,
			LatestBadAnswer:          null.TimeFrom(time.Now()),
			LatestBadAnswerWasToday:  true,
			GoodAnswers:              0,
			GoodAnswersToday:         0,
			LatestGoodAnswer:         null.Time{},
			LatestGoodAnswerWasToday: false,
		},
	}

	expectedExercises := []frontend.Exercise{
		{
			Id:                       1,
			Question:                 "question",
			Answer:                   "answer",
			BadAnswers:               2,
			BadAnswersToday:          1,
			LatestBadAnswer:          null.TimeFrom(time.Now()),
			LatestBadAnswerWasToday:  true,
			GoodAnswers:              0,
			GoodAnswersToday:         0,
			LatestGoodAnswer:         null.Time{},
			LatestGoodAnswerWasToday: false,
		},
	}

	responseBody, err := json.Marshal(exercises)
	suite.Assert().NoError(err)

	method := "GET"
	route := server.FetchExercises
	params := map[string]string{"lesson_id": "1"}
	reqBody := []byte(nil)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return(responseBody)

	result, err := suite.client.FetchExercises(ctx, frontend.Lesson{Id: lesson.Id})

	suite.Nil(err)
	suite.Equal(expectedExercises[0].Id, result[0].Id)
	suite.Equal(expectedExercises[0].Question, result[0].Question)
	suite.Equal(expectedExercises[0].Answer, result[0].Answer)
	suite.Equal(expectedExercises[0].BadAnswers, result[0].BadAnswers)
	suite.Equal(expectedExercises[0].BadAnswersToday, result[0].BadAnswersToday)
	suite.Equal(expectedExercises[0].LatestBadAnswer.Time.Local().Format("Mon Jan 2 15:04:05"), result[0].LatestBadAnswer.Time.Local().Format("Mon Jan 2 15:04:05"))
	suite.Equal(expectedExercises[0].LatestBadAnswerWasToday, result[0].LatestBadAnswerWasToday)
	suite.Equal(expectedExercises[0].GoodAnswers, result[0].GoodAnswers)
	suite.Equal(expectedExercises[0].GoodAnswersToday, result[0].GoodAnswersToday)
	suite.Equal(expectedExercises[0].LatestGoodAnswer.Time.Local().Format("Mon Jan 2 15:04:05"), result[0].LatestGoodAnswer.Time.Local().Format("Mon Jan 2 15:04:05"))
	suite.Equal(expectedExercises[0].LatestGoodAnswerWasToday, result[0].LatestGoodAnswerWasToday)
	suite.caller.AssertExpectations(suite.T())
}

func (suite *ClientSuite) TestClient_UpsertLesson() {
	ctx := context.Background()

	lesson := frontend.Lesson{}

	method := "POST"
	route := server.UpsertLesson
	params := map[string]string(nil)
	reqBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""))

	err = suite.client.UpsertLesson(ctx, lesson)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_DeleteLesson() {
	ctx := context.Background()

	lesson := frontend.Lesson{}

	method := "POST"
	route := server.DeleteLesson
	params := map[string]string(nil)
	reqBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""))

	err = suite.client.DeleteLesson(ctx, lesson)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_UpsertExercise() {
	ctx := context.Background()

	exercise := frontend.Exercise{}

	method := "POST"
	route := server.UpsertExercise
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""))

	err = suite.client.UpsertExercise(ctx, exercise)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_StoreExercises() {
	ctx := context.Background()

	var exercises []frontend.Exercise

	method := "POST"
	route := server.StoreExercises
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercises)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""))

	err = suite.client.StoreExercises(ctx, exercises)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_DeleteExercise() {
	ctx := context.Background()

	exercise := frontend.Exercise{}

	method := "POST"
	route := server.DeleteExercise
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""))

	err = suite.client.DeleteExercise(ctx, exercise)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_StoreResult() {
	ctx := context.Background()

	result := frontend.Result{}

	method := "POST"
	route := server.StoreResult
	params := map[string]string(nil)
	reqBody, err := json.Marshal(result)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""))

	err = suite.client.StoreResult(ctx, result)
	suite.Assert().NoError(err)
}
