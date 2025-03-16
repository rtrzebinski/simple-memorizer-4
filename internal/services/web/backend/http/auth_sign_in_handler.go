package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type AuthSignInHandler struct {
	s Service
}

func NewAuthSignInHandler(s Service) *AuthSignInHandler {
	return &AuthSignInHandler{s: s}
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

	accessToken, err := h.s.SignIn(ctx, signInRequest.Email, signInRequest.Password)

	if err != nil {
		log.Print(fmt.Errorf("failed to sign in user: %w", err))
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	signInResponse := backend.SignInResponse{
		AccessToken: accessToken,
	}

	encoded, err := json.Marshal(signInResponse)

	if err != nil {
		log.Print(fmt.Errorf("failed to encode AuthSignInHandler HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = res.Write(encoded)
}
