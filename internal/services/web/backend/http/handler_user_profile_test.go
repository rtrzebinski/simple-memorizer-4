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

func TestHandlerUserProfile_success(t *testing.T) {
	v := NewTokenVerifierMock()
	v.On("VerifyAndUser", mock.Anything, "accessToken").Return(&backend.User{ID: "userID", Name: "Test Name", Email: "test@example.com"}, nil)
	r := NewTokenRefresherMock()
	r.On("Refresh", mock.Anything, mock.Anything).Return(backend.Tokens{}, nil).Maybe()

	route := auth(v, r, false)(NewHandlerUserProfile(v))

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, UserProfile, nil)
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "access_token", Value: "accessToken"})

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var profile backend.UserProfile
	err = json.Unmarshal(res.Body.Bytes(), &profile)
	assert.NoError(t, err)
	assert.Equal(t, "Test Name", profile.Name)
	assert.Equal(t, "test@example.com", profile.Email)

	v.AssertExpectations(t)
	r.AssertExpectations(t)
}

func TestHandlerUserProfile_unauthorized(t *testing.T) {
	v := NewTokenVerifierMock()
	r := NewTokenRefresherMock()

	route := auth(v, r, false)(NewHandlerUserProfile(v))

	res := httptest.NewRecorder()
	req, err := http.NewRequest(http.MethodGet, UserProfile, nil)
	assert.NoError(t, err)

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusUnauthorized, res.Code)
}
