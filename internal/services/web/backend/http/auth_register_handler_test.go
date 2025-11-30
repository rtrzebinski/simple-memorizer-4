package http

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/stretchr/testify/assert"
)

func TestAuthRegisterHandler(t *testing.T) {
	ctx := context.Background()

	input := backend.RegisterRequest{
		FirstName: "firstname",
		LastName:  "lastname",
		Email:     "email",
		Password:  "password",
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()
	service.On("Register", ctx, input.FirstName, input.LastName, input.Email, input.Password).Return("accessToken", nil)

	handler := NewAuthRegisterHandler(service)

	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	service.AssertExpectations(t)
	var registerResponse backend.RegisterResponse
	err = json.Unmarshal(res.Body.Bytes(), &registerResponse)
	assert.NoError(t, err)
	assert.Equal(t, "accessToken", registerResponse.AccessToken)
}

func TestAuthRegisterHandler_unauthorized(t *testing.T) {
	ctx := context.Background()

	input := backend.RegisterRequest{
		FirstName: "firstname",
		LastName:  "lastname",
		Email:     "email",
		Password:  "password",
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()
	service.On("Register", ctx, input.FirstName, input.LastName, input.Email, input.Password).Return("", errors.New("unauthorized"))

	handler := NewAuthRegisterHandler(service)

	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Code)
	service.AssertExpectations(t)
}
