package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/auth"
)

type UserProfileHandler struct {
	v TokenVerifier
}

func NewUserProfileHandler(v TokenVerifier) *UserProfileHandler {
	return &UserProfileHandler{
		v: v,
	}
}

func (h *UserProfileHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
		log.Print(fmt.Errorf("failed to encode UserProfileHandler HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write UserProfileHandler HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
