package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/guregu/null/v5"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchExercisesOfLessonHandler(t *testing.T) {
	exercise := backend.Exercise{
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
	}
	exercises := backend.Exercises{exercise}

	lessonId := 10

	service := NewServiceMock()

	oldestExerciseID := 1

	service.On("FetchExercises", mock.Anything, backend.Lesson{Id: lessonId}, oldestExerciseID, "100").Return(exercises, nil)

	v := NewTokenVerifierMock()
	v.On("VerifyAndUserID", mock.Anything, mock.Anything).Return("100", nil)
	route := RequireAuth(v)(NewFetchExercisesOfLessonHandler(service))

	u, _ := url.Parse("/")
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lessonId))
	params.Add("oldest_exercise_id", strconv.Itoa(oldestExerciseID))
	u.RawQuery = params.Encode()

	req := &http.Request{}
	req.URL = u
	req.Header = make(map[string][]string)
	req.Header.Set("authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMDAifQ.bEOa2kaRwC1f7Ow-7WgSltYq-Vz9JUDCo3EPe7KEXd8")

	res := httptest.NewRecorder()

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var result backend.Exercises
	json.Unmarshal(res.Body.Bytes(), &result)

	assert.Equal(t, exercises[0].Id, result[0].Id)
	assert.Equal(t, exercises[0].Question, result[0].Question)
	assert.Equal(t, exercises[0].Answer, result[0].Answer)
	assert.Equal(t, exercises[0].BadAnswers, result[0].BadAnswers)
	assert.Equal(t, exercises[0].BadAnswersToday, result[0].BadAnswersToday)
	assert.Equal(t, exercises[0].LatestBadAnswer.Time.Format("Mon Jan 2 15:04:05"), result[0].LatestBadAnswer.Time.Format("Mon Jan 2 15:04:05"))
	assert.Equal(t, exercises[0].LatestBadAnswerWasToday, result[0].LatestBadAnswerWasToday)
	assert.Equal(t, exercises[0].GoodAnswers, result[0].GoodAnswers)
	assert.Equal(t, exercises[0].GoodAnswersToday, result[0].GoodAnswersToday)
	assert.Equal(t, exercises[0].LatestGoodAnswer.Time.Format("Mon Jan 2 15:04:05"), result[0].LatestGoodAnswer.Time.Format("Mon Jan 2 15:04:05"))
	assert.Equal(t, exercises[0].LatestGoodAnswerWasToday, result[0].LatestGoodAnswerWasToday)
}

func TestFetchExercisesOfLessonHandler_invalidInput(t *testing.T) {
	service := NewServiceMock()

	v := NewTokenVerifierMock()
	v.On("VerifyAndUserID", mock.Anything, mock.Anything).Return("100", nil)
	route := RequireAuth(v)(NewFetchExercisesOfLessonHandler(service))

	u, _ := url.Parse("/")

	req := &http.Request{}
	req.URL = u
	req.Header = make(map[string][]string)
	req.Header.Set("authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMDAifQ.bEOa2kaRwC1f7Ow-7WgSltYq-Vz9JUDCo3EPe7KEXd8")

	res := httptest.NewRecorder()

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err := json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateLessonIdentified(backend.Lesson{}).Error(), result)
}
