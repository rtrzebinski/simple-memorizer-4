package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/server/validation"
	"github.com/stretchr/testify/assert"
)

func TestFetchExercisesOfLessonHandler(t *testing.T) {
	exercise := backend.Exercise{
		Id:       1,
		Question: "question",
		Answer:   "answer",
	}
	exercises := backend.Exercises{exercise}

	lessonId := 10

	reader := NewReaderMock()
	reader.On("FetchExercises", backend.Lesson{Id: lessonId}).Return(exercises)

	route := NewFetchExercisesOfLessonHandler(reader)

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

	assert.Equal(t, exercises, result)
}

func TestFetchExercisesOfLessonHandler_invalidInput(t *testing.T) {
	reader := NewReaderMock()

	route := NewFetchExercisesOfLessonHandler(reader)

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
