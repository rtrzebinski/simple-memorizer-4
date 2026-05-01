package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerFetchLessons(t *testing.T) {
	lesson := backend.Lesson{}
	lessons := backend.Lessons{lesson}

	service := NewServiceMock()
	service.On("FetchLessons", mock.Anything, "100").Return(lessons, nil)

	v := NewTokenVerifierMock()
	v.On("VerifyAndUser", mock.Anything, "accessToken").Return(&backend.User{ID: "100"}, nil)
	r := NewTokenRefresherMock()
	route := Auth(v, r, false)(NewHandlerFetchLessons(service))

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, FetchLessons, nil)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "accessToken"})

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var result backend.Lessons
	err := json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, lessons, result)

	service.AssertExpectations(t)
	v.AssertExpectations(t)
	r.AssertExpectations(t)
}

func TestHandlerFetchLessons_unauthorized(t *testing.T) {
	service := NewServiceMock()

	v := NewTokenVerifierMock()
	r := NewTokenRefresherMock()
	route := Auth(v, r, false)(NewHandlerFetchLessons(service))

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, FetchLessons, nil)

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Code)
}
