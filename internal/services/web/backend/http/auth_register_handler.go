package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type AuthRegisterHandler struct {
	s Service
}

func NewAuthRegisterHandler(s Service) *AuthRegisterHandler {
	return &AuthRegisterHandler{s: s}
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

	accessToken, err := h.s.Register(ctx, registerRequest.Name, registerRequest.Email, registerRequest.Password)

	if err != nil {
		log.Print(fmt.Errorf("failed to register user: %w", err))
		res.WriteHeader(http.StatusUnauthorized)
		return
	}

	registerResponse := backend.RegisterResponse{
		AccessToken: accessToken,
	}

	encoded, err := json.Marshal(registerResponse)

	if err != nil {
		log.Print(fmt.Errorf("failed to encode AuthRegisterHandler HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = res.Write(encoded)
}
