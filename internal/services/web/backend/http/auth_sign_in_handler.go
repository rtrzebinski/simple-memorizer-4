package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type AuthSignInHandler struct {
	s      Service
	secure bool
}

func NewAuthSignInHandler(s Service, secure bool) *AuthSignInHandler {
	return &AuthSignInHandler{
		s:      s,
		secure: secure,
	}
}

func (h *AuthSignInHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var signInRequest backend.SignInRequest

	err := json.NewDecoder(req.Body).Decode(&signInRequest)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode AuthSignInHandler HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	t, err := h.s.SignIn(ctx, signInRequest.Email, signInRequest.Password)
	if err != nil {
		log.Print(fmt.Errorf("failed to sign in user: %w", err))
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
