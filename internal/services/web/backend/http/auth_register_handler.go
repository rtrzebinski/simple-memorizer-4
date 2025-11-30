package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type AuthRegisterHandler struct {
	s      Service
	secure bool
}

func NewAuthRegisterHandler(s Service, secure bool) *AuthRegisterHandler {
	return &AuthRegisterHandler{
		s:      s,
		secure: secure,
	}
}

func (h *AuthRegisterHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var registerRequest backend.RegisterRequest

	err := json.NewDecoder(req.Body).Decode(&registerRequest)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode AuthRegisterHandler HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	t, err := h.s.Register(ctx, registerRequest.FirstName, registerRequest.LastName, registerRequest.Email, registerRequest.Password)

	if err != nil {
		log.Print(fmt.Errorf("failed to register user: %w", err))
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	// set refresh token in HttpOnly cookie
	http.SetCookie(res, &http.Cookie{
		Name:     "refresh_token",
		Value:    t.RefreshToken,
		Path:     "/",
		MaxAge:   t.RefreshExpiresIn,
		Secure:   h.secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	// set access token in HttpOnly cookie
	http.SetCookie(res, &http.Cookie{
		Name:     "access_token",
		Value:    t.AccessToken,
		Path:     "/",
		MaxAge:   t.ExpiresIn,
		Secure:   h.secure,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})
}
