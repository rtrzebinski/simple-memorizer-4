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
	"github.com/stretchr/testify/mock"
)

func TestExportLessonCsvHandler(t *testing.T) {
	exercise := backend.Exercise{
		Id:       1,
		Question: "question",
		Answer:   "answer",
	}
	exercises := backend.Exercises{exercise}

	lesson := backend.Lesson{Id: 2}

	reader := NewReaderMock()
	reader.On("FetchExercises", lesson).Return(exercises)
	reader.On("HydrateLesson", &lesson).Run(func(args mock.Arguments) {
		args.Get(0).(*backend.Lesson).Name = "lesson name"
	})

	route := NewExportLessonCsvHandler(reader)

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

func TestExportLessonCsvHandler_invalidInput(t *testing.T) {
	reader := NewReaderMock()

	route := NewExportLessonCsvHandler(reader)

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
