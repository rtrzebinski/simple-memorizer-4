package server

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
)

func TestExportLessonCsv(t *testing.T) {
	exercise := models.Exercise{
		Id:       1,
		Question: "question",
		Answer:   "answer",
	}
	exercises := models.Exercises{exercise}

	lesson := models.Lesson{Id: 2}

	reader := storage.NewReaderMock()
	reader.On("FetchExercises", lesson).Return(exercises)
	reader.On("HydrateLesson", &lesson).Run(func(args mock.Arguments) {
		args.Get(0).(*models.Lesson).Name = "lesson name"
	})

	route := NewExportLessonCsv(reader)

	u, _ := url.Parse("/")
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lesson.Id))
	u.RawQuery = params.Encode()

	req := &http.Request{}
	req.URL = u

	res := httptest.NewRecorder()

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "question,answer\n", string(res.Body.Bytes()))
	assert.Equal(t, "attachment; filename=lesson name.csv", res.Header().Get("Content-Disposition"))
	assert.Equal(t, "application/octet-stream", res.Header().Get("Content-Type"))
	assert.Equal(t, "16", res.Header().Get("Content-Length"))
}

func TestExportLessonCsv_invalidInput(t *testing.T) {
	reader := storage.NewReaderMock()

	route := NewExportLessonCsv(reader)

	u, _ := url.Parse("/")

	req := &http.Request{}
	req.URL = u

	res := httptest.NewRecorder()

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err := json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateLessonIdentified(models.Lesson{}).Error(), result)
}
