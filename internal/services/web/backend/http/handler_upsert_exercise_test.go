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

func TestHandlerUpsertExercise(t *testing.T) {
	input := backend.Exercise{
		Question: "question",
		Answer:   "answer",
		Lesson:   &backend.Lesson{Id: 10},
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()
	service.On("UpsertExercise", mock.Anything, &input, "100").Return(nil)

	v := NewTokenVerifierMock()
	v.On("VerifyAndUser", mock.Anything, "accessToken").Return(&backend.User{ID: "100"}, nil)
	r := NewTokenRefresherMock()
	route := Auth(v, r, false)(NewHandlerUpsertExercise(service))

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, UpsertExercise, strings.NewReader(string(body)))
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "accessToken"})

	route.ServeHTTP(res, req)

	service.AssertExpectations(t)
	v.AssertExpectations(t)
	r.AssertExpectations(t)
}

func TestHandlerUpsertExercise_invalidInput(t *testing.T) {
	input := backend.Exercise{}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()

	v := NewTokenVerifierMock()
	v.On("VerifyAndUser", mock.Anything, "accessToken").Return(&backend.User{ID: "100"}, nil)
	r := NewTokenRefresherMock()
	route := Auth(v, r, false)(NewHandlerUpsertExercise(service))

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, UpsertExercise, strings.NewReader(string(body)))
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "accessToken"})

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateUpsertExercise(input, nil).Error(), result)

	service.AssertExpectations(t)
	v.AssertExpectations(t)
	r.AssertExpectations(t)
}

func TestHandlerUpsertExercise_unauthorized(t *testing.T) {
	service := NewServiceMock()

	v := NewTokenVerifierMock()
	r := NewTokenRefresherMock()
	route := Auth(v, r, false)(NewHandlerUpsertExercise(service))

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, UpsertExercise, strings.NewReader(`{}`))

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Code)
}
