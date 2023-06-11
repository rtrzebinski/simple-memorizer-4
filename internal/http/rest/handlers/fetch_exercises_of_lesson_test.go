package rest

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/validators"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestFetchExercisesOfLesson(t *testing.T) {
	exercise := models.Exercise{
		Id:          1,
		Question:    "question",
		Answer:      "answer",
		BadAnswers:  2,
		GoodAnswers: 3,
	}
	exercises := models.Exercises{exercise}

	lessonId := 10

	reader := storage.NewReaderMock()
	reader.On("FetchExercisesOfLesson", models.Lesson{Id: lessonId}).Return(exercises)

	route := NewFetchExercisesOfLesson(reader)

	u, _ := url.Parse("/")
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lessonId))
	u.RawQuery = params.Encode()

	req := &http.Request{}
	req.URL = u

	res := httptest.NewRecorder()

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var result models.Exercises
	json.Unmarshal(res.Body.Bytes(), &result)

	assert.Equal(t, exercises, result)
}

func TestFetchExercisesOfLesson_invalidInput(t *testing.T) {
	reader := storage.NewReaderMock()

	route := NewFetchExercisesOfLesson(reader)

	u, _ := url.Parse("/")

	req := &http.Request{}
	req.URL = u

	res := httptest.NewRecorder()

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err := json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validators.ValidateLessonIdentified(models.Lesson{}).Error(), result)
}
