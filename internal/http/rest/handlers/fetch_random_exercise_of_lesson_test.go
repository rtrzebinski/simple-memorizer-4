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

func TestFetchRandomExerciseOfLesson(t *testing.T) {
	exercise := models.Exercise{
		Id:          1,
		Question:    "question",
		Answer:      "answer",
		BadAnswers:  2,
		GoodAnswers: 3,
	}

	lessonId := 10

	reader := storage.NewReaderMock()
	reader.On("FetchRandomExerciseOfLesson", models.Lesson{Id: lessonId}).Return(exercise)

	route := NewFetchRandomExerciseOfLesson(reader)

	u, _ := url.Parse("/")
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lessonId))
	u.RawQuery = params.Encode()

	req := &http.Request{}
	req.URL = u

	res := httptest.NewRecorder()

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var result models.Exercise
	json.Unmarshal(res.Body.Bytes(), &result)

	assert.Equal(t, exercise, result)
}

func TestFetchRandomExerciseOfLesson_invalidInput(t *testing.T) {
	reader := storage.NewReaderMock()

	route := NewFetchRandomExerciseOfLesson(reader)

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
