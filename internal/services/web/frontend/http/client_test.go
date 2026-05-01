package http

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/guregu/null/v5"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http"
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
	ctx := app.Context{Context: suite.T().Context()}

	lessons := []frontend.Lesson{{Name: "name"}}

	responseBody, err := json.Marshal(lessons)
	suite.Assert().NoError(err)

	method := "GET"
	route := http.FetchLessons
	params := map[string]string(nil)
	reqBody := []byte(nil)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return(responseBody)

	result, err := suite.client.FetchLessons(ctx)
	suite.Assert().NoError(err)
	suite.Assert().Equal(lessons, result)
}

func (suite *ClientSuite) TestClient_HydrateLesson() {
	ctx := app.Context{Context: suite.T().Context()}

	lesson := &frontend.Lesson{Id: 10}

	responseBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	method := "GET"
	route := http.HydrateLesson
	params := map[string]string{"lesson_id": "10"}
	reqBody := []byte(nil)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return(responseBody)

	err = suite.client.HydrateLesson(ctx, lesson)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_FetchExercises() {
	ctx := app.Context{Context: suite.T().Context()}

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

	oldestExerciseID := 1

	method := "GET"
	route := http.FetchExercises
	params := map[string]string{"lesson_id": "1", "oldest_exercise_id": "1"}
	reqBody := []byte(nil)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return(responseBody)

	result, err := suite.client.FetchExercises(ctx, frontend.Lesson{Id: lesson.Id}, oldestExerciseID)

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
	ctx := app.Context{Context: suite.T().Context()}

	lesson := frontend.Lesson{}

	method := "POST"
	route := http.UpsertLesson
	params := map[string]string(nil)
	reqBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""))

	err = suite.client.UpsertLesson(ctx, lesson)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_DeleteLesson() {
	ctx := app.Context{Context: suite.T().Context()}

	lesson := frontend.Lesson{}

	method := "POST"
	route := http.DeleteLesson
	params := map[string]string(nil)
	reqBody, err := json.Marshal(lesson)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""))

	err = suite.client.DeleteLesson(ctx, lesson)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_UpsertExercise() {
	ctx := app.Context{Context: suite.T().Context()}

	exercise := frontend.Exercise{}

	method := "POST"
	route := http.UpsertExercise
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""))

	err = suite.client.UpsertExercise(ctx, exercise)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_StoreExercises() {
	ctx := app.Context{Context: suite.T().Context()}

	var exercises []frontend.Exercise

	method := "POST"
	route := http.StoreExercises
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercises)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""))

	err = suite.client.StoreExercises(ctx, exercises)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_DeleteExercise() {
	ctx := app.Context{Context: suite.T().Context()}

	exercise := frontend.Exercise{}

	method := "POST"
	route := http.DeleteExercise
	params := map[string]string(nil)
	reqBody, err := json.Marshal(exercise)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""))

	err = suite.client.DeleteExercise(ctx, exercise)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_StoreResult() {
	ctx := app.Context{Context: suite.T().Context()}

	result := frontend.Result{}

	method := "POST"
	route := http.StoreResult
	params := map[string]string(nil)
	reqBody, err := json.Marshal(result)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""))

	err = suite.client.StoreResult(ctx, result)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_AuthRegister() {
	ctx := app.Context{Context: suite.T().Context()}

	req := frontend.RegisterRequest{}

	method := "POST"
	route := http.AuthRegister
	params := map[string]string(nil)
	reqBody, err := json.Marshal(req)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""), nil)

	err = suite.client.AuthRegister(ctx, req)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_AuthSignIn() {
	ctx := app.Context{Context: suite.T().Context()}

	req := frontend.SignInRequest{}

	method := "POST"
	route := http.AuthSignIn
	params := map[string]string(nil)
	reqBody, err := json.Marshal(req)
	suite.Assert().NoError(err)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""), nil)

	err = suite.client.AuthSignIn(ctx, req)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_AuthLogout() {
	ctx := app.Context{Context: suite.T().Context()}

	method := "POST"
	route := http.AuthLogout
	params := map[string]string(nil)
	reqBody := []byte(nil)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return([]byte(""), nil)

	err := suite.client.AuthLogout(ctx)
	suite.Assert().NoError(err)
}

func (suite *ClientSuite) TestClient_UserProfile() {
	ctx := app.Context{Context: suite.T().Context()}

	user := &frontend.UserProfile{
		Name:  "John Doe",
		Email: "john@example.com",
	}

	responseBody, err := json.Marshal(user)
	suite.Assert().NoError(err)

	method := "GET"
	route := http.UserProfile
	params := map[string]string(nil)
	reqBody := []byte(nil)

	suite.caller.On("Call", ctx, method, route, params, reqBody).Return(responseBody, nil)

	result, err := suite.client.UserProfile(ctx)
	suite.Assert().NoError(err)
	suite.Assert().Equal(user, result)
}
