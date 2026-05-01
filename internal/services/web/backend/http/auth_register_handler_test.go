package http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthRegisterHandler(t *testing.T) {
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

	tokens := backend.Tokens{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
	}

	service := NewServiceMock()
	service.On("Register", mock.Anything, input.FirstName, input.LastName, input.Email, input.Password).Return(tokens, nil)

	handler := NewAuthRegisterHandler(service, true)

	req, err := http.NewRequest(http.MethodPost, AuthRegister, io.NopCloser(strings.NewReader(string(body))))
	assert.NoError(t, err)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	service.AssertExpectations(t)

	refreshTokenCookie := res.Result().Cookies()[0]
	assert.Equal(t, "refresh_token", refreshTokenCookie.Name)
	assert.Equal(t, "refreshToken", refreshTokenCookie.Value)
	assert.True(t, refreshTokenCookie.HttpOnly)
	assert.True(t, refreshTokenCookie.Secure)
	assert.Equal(t, http.SameSiteStrictMode, refreshTokenCookie.SameSite)

	accessTokenCookie := res.Result().Cookies()[1]
	assert.Equal(t, "access_token", accessTokenCookie.Name)
	assert.Equal(t, "accessToken", accessTokenCookie.Value)
	assert.True(t, accessTokenCookie.HttpOnly)
	assert.True(t, accessTokenCookie.Secure)
	assert.Equal(t, http.SameSiteStrictMode, accessTokenCookie.SameSite)
}

func TestAuthRegisterHandler_unauthorized(t *testing.T) {
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
	service.On("Register", mock.Anything, input.FirstName, input.LastName, input.Email, input.Password).Return(backend.Tokens{}, errors.New("unauthorized"))

	handler := NewAuthRegisterHandler(service, true)

	req, err := http.NewRequest(http.MethodPost, AuthRegister, io.NopCloser(strings.NewReader(string(body))))
	assert.NoError(t, err)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Code)
	service.AssertExpectations(t)
}
