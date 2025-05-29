package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestExportLessonCsvHandler(t *testing.T) {
	ctx := context.Background()

	exercise := backend.Exercise{
		Id:       1,
		Question: "question",
		Answer:   "answer",
	}
	exercises := backend.Exercises{exercise}

	lesson := backend.Lesson{Id: 2}

	service := NewServiceMock()

	oldestExerciseID := 1

	service.On("FetchExercises", ctx, lesson, oldestExerciseID, "").Return(exercises, nil)
	service.On("HydrateLesson", ctx, &lesson, "").Run(func(args mock.Arguments) {
		args.Get(1).(*backend.Lesson).Name = "lesson name"
	}).Return(nil)

	route := NewExportLessonCsvHandler(service)

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
	service := NewServiceMock()

	route := NewExportLessonCsvHandler(service)

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
