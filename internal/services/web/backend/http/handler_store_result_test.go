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

func TestHandlerStoreResult_goodAnswer(t *testing.T) {
	input := backend.Result{
		Exercise: &backend.Exercise{
			Id: 10,
		},
		Type: backend.Good,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()
	service.On("PublishGoodAnswer", mock.Anything, "userID", input.Exercise.Id).Return(nil)

	v := NewTokenVerifierMock()
	v.On("VerifyAndUser", mock.Anything, "accessToken").Return(&backend.User{ID: "userID"}, nil)
	r := NewTokenRefresherMock()
	route := auth(v, r, false)(NewHandlerStoreResult(service))

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, StoreResult, strings.NewReader(string(body)))
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "accessToken"})

	route.ServeHTTP(res, req)

	service.AssertExpectations(t)
	v.AssertExpectations(t)
	r.AssertExpectations(t)
}

func TestHandlerStoreResult_badAnswer(t *testing.T) {
	input := backend.Result{
		Exercise: &backend.Exercise{
			Id: 10,
		},
		Type: backend.Bad,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()
	service.On("PublishBadAnswer", mock.Anything, "userID", input.Exercise.Id).Return(nil)

	v := NewTokenVerifierMock()
	v.On("VerifyAndUser", mock.Anything, "accessToken").Return(&backend.User{ID: "userID"}, nil)
	r := NewTokenRefresherMock()
	route := auth(v, r, false)(NewHandlerStoreResult(service))

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, StoreResult, strings.NewReader(string(body)))
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "accessToken"})

	route.ServeHTTP(res, req)

	service.AssertExpectations(t)
	v.AssertExpectations(t)
	r.AssertExpectations(t)
}

func TestHandlerStoreResult_invalidInput(t *testing.T) {
	input := backend.Result{}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()

	v := NewTokenVerifierMock()
	v.On("VerifyAndUser", mock.Anything, "accessToken").Return(&backend.User{ID: "userID"}, nil)
	r := NewTokenRefresherMock()
	route := auth(v, r, false)(NewHandlerStoreResult(service))

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, StoreResult, strings.NewReader(string(body)))
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "accessToken"})

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateStoreResult(input).Error(), result)

	service.AssertExpectations(t)
	v.AssertExpectations(t)
	r.AssertExpectations(t)
}

func TestHandlerStoreResult_unauthorized(t *testing.T) {
	service := NewServiceMock()

	v := NewTokenVerifierMock()
	r := NewTokenRefresherMock()
	route := auth(v, r, false)(NewHandlerStoreResult(service))

	res := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, StoreResult, strings.NewReader(`{}`))

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Code)
}
