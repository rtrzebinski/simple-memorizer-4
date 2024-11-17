package server

import (
	"encoding/json"
	"github.com/guregu/null/v5"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/server/validation"
	"github.com/stretchr/testify/assert"
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

	service.On("FetchExercises", backend.Lesson{Id: lessonId}).Return(exercises, nil)

	route := NewFetchExercisesOfLessonHandler(service)

	u, _ := url.Parse("/")
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lessonId))
	u.RawQuery = params.Encode()

	req := &http.Request{}
	req.URL = u

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

	route := NewFetchExercisesOfLessonHandler(service)

	u, _ := url.Parse("/")

	req := &http.Request{}
	req.URL = u

	res := httptest.NewRecorder()

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err := json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateLessonIdentified(backend.Lesson{}).Error(), result)
}
