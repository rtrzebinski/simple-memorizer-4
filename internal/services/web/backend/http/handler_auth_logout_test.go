package http

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestHandlerAuthLogout(t *testing.T) {
	service := NewServiceMock()
	service.On("Revoke", mock.Anything, "refreshToken").Return(nil)

	handler := NewHandlerAuthLogout(service)

	req, err := http.NewRequest(http.MethodPost, AuthLogout, nil)
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refreshToken"})
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
	service.AssertExpectations(t)

	refreshTokenCookie := res.Result().Cookies()[0]
	assert.Equal(t, "refresh_token", refreshTokenCookie.Name)
	assert.Equal(t, "", refreshTokenCookie.Value)
	assert.True(t, refreshTokenCookie.HttpOnly)
	assert.Equal(t, -1, refreshTokenCookie.MaxAge)

	accessTokenCookie := res.Result().Cookies()[1]
	assert.Equal(t, "access_token", accessTokenCookie.Name)
	assert.Equal(t, "", accessTokenCookie.Value)
	assert.True(t, accessTokenCookie.HttpOnly)
	assert.Equal(t, -1, accessTokenCookie.MaxAge)
}

func TestHandlerAuthLogout_noRefreshCookie(t *testing.T) {
	handler := NewHandlerAuthLogout(NewServiceMock())

	req, err := http.NewRequest(http.MethodPost, AuthLogout, nil)
	assert.NoError(t, err)
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)
}

func TestHandlerAuthLogout_revokeFailed(t *testing.T) {
	service := NewServiceMock()
	service.On("Revoke", mock.Anything, "refreshToken").Return(errors.New("revoke failed"))

	handler := NewHandlerAuthLogout(service)

	req, err := http.NewRequest(http.MethodPost, AuthLogout, nil)
	assert.NoError(t, err)
	req.AddCookie(&http.Cookie{Name: "refresh_token", Value: "refreshToken"})
	res := httptest.NewRecorder()

	handler.ServeHTTP(res, req)

	assert.Equal(t, http.StatusInternalServerError, res.Code)
	service.AssertExpectations(t)
}
