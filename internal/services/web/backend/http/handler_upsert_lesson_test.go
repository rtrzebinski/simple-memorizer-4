package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerUpsertLesson(t *testing.T) {
	input := backend.Lesson{
		Name:        "name",
		Description: "description",
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()
	service.On("UpsertLesson", mock.Anything, "userID", &input).Return(nil)

	v := NewTokenVerifierMock()
	v.On("VerifyAndUser", mock.Anything, "accessToken").Return(&backend.User{ID: "userID"}, nil)
	r := NewTokenRefresherMock()
	route := auth(v, r, false)(NewHandlerUpsertLesson(service))

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, UpsertLesson, strings.NewReader(string(body)))
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "accessToken"})

	route.ServeHTTP(res, req)

	service.AssertExpectations(t)
	v.AssertExpectations(t)
	r.AssertExpectations(t)
}

func TestHandlerUpsertLesson_invalidInput(t *testing.T) {
	input := backend.Lesson{}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()

	v := NewTokenVerifierMock()
	v.On("VerifyAndUser", mock.Anything, "accessToken").Return(&backend.User{ID: "userID"}, nil)
	r := NewTokenRefresherMock()
	route := auth(v, r, false)(NewHandlerUpsertLesson(service))

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, UpsertLesson, strings.NewReader(string(body)))
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "accessToken"})

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateUpsertLesson(input, nil).Error(), result)

	service.AssertExpectations(t)
	v.AssertExpectations(t)
	r.AssertExpectations(t)
}

func TestHandlerUpsertLesson_unauthorized(t *testing.T) {
	service := NewServiceMock()

	v := NewTokenVerifierMock()
	r := NewTokenRefresherMock()
	route := auth(v, r, false)(NewHandlerUpsertLesson(service))

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, UpsertLesson, strings.NewReader(`{}`))

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Code)
}
