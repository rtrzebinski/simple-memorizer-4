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

func TestNewAuthSignInHandler(t *testing.T) {
	ctx := context.Background()

	input := backend.SignInRequest{
		Email:    "email",
		Password: "password",
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()
	service.On("SignIn", ctx, input.Email, input.Password).Return("accessToken", nil)

	handler := NewAuthSignInHandler(service)

	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	service.AssertExpectations(t)
	var signInResponse backend.SignInResponse
	err = json.Unmarshal(res.Body.Bytes(), &signInResponse)
	assert.NoError(t, err)
	assert.Equal(t, "accessToken", signInResponse.AccessToken)
}

func TestNewAuthSignInHandler_unauthorized(t *testing.T) {
	ctx := context.Background()

	input := backend.SignInRequest{
		Email:    "email",
		Password: "password",
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()
	service.On("SignIn", ctx, input.Email, input.Password).Return("", errors.New("unauthorized"))

	handler := NewAuthSignInHandler(service)

	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Code)
	service.AssertExpectations(t)
}
