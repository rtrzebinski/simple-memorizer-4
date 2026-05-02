package http

import (
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

func TestHandlerExportLessonCsv(t *testing.T) {
	exercise1 := backend.Exercise{
		Id:       1,
		Question: "question1",
		Answer:   "answer1",
	}
	exercise2 := backend.Exercise{
		Id:       2,
		Question: "question2",
		Answer:   "answer2",
	}
	exercises := backend.Exercises{exercise2, exercise1}

	lesson := backend.Lesson{Id: 2}

	service := NewServiceMock()

	oldestExerciseID := 1

	service.On("FetchExercises", mock.Anything, "userID", lesson, oldestExerciseID).Return(exercises, nil)
	service.On("HydrateLesson", mock.Anything, "userID", &lesson).Run(func(args mock.Arguments) {
		args.Get(2).(*backend.Lesson).Name = "lesson name"
	}).Return(nil)

	v := NewTokenVerifierMock()
	v.On("VerifyAndUser", mock.Anything, "accessToken").Return(&backend.User{ID: "userID"}, nil)
	r := NewTokenRefresherMock()
	route := Auth(v, r, false)(NewHandlerExportLessonCsv(service))

	u, _ := url.Parse(ExportLessonCsv)
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lesson.Id))
	u.RawQuery = params.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "accessToken"})

	res := httptest.NewRecorder()

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	assert.Equal(t, "question1,answer1\nquestion2,answer2\n", string(res.Body.Bytes()))
	assert.Equal(t, "attachment; filename=lesson name.csv", res.Header().Get("Content-Disposition"))
	assert.Equal(t, "application/octet-stream", res.Header().Get("Content-Type"))
	assert.Equal(t, "36", res.Header().Get("Content-Length"))
}

func TestHandlerExportLessonCsv_invalidInput(t *testing.T) {
	service := NewServiceMock()

	v := NewTokenVerifierMock()
	v.On("VerifyAndUser", mock.Anything, "accessToken").Return(&backend.User{ID: "userID"}, nil)
	r := NewTokenRefresherMock()
	route := Auth(v, r, false)(NewHandlerExportLessonCsv(service))

	req, err := http.NewRequest(http.MethodGet, ExportLessonCsv, nil)
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "accessToken"})

	res := httptest.NewRecorder()

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateLessonIdentified(backend.Lesson{}).Error(), result)
}

func TestHandlerExportLessonCsv_unauthorized(t *testing.T) {
	service := NewServiceMock()

	v := NewTokenVerifierMock()
	r := NewTokenRefresherMock()
	route := Auth(v, r, false)(NewHandlerExportLessonCsv(service))

	req, _ := http.NewRequest(http.MethodGet, ExportLessonCsv, nil)

	res := httptest.NewRecorder()

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Code)
}
