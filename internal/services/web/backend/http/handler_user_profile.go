package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/auth"
)

type HandlerUserProfile struct {
	v TokenVerifier
}

func NewHandlerUserProfile(v TokenVerifier) *HandlerUserProfile {
	return &HandlerUserProfile{
		v: v,
	}
}

func (h *HandlerUserProfile) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	profile, ok := auth.UserProfileFromContext(req.Context())
	if !ok || profile == nil {
		http.Error(res, "unauthorized", http.StatusUnauthorized)
		return
	}

	encoded, err := json.Marshal(profile)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode HandlerUserProfile HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write HandlerUserProfile HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
