package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/auth"
)

type FetchLessonsHandler struct {
	s Service
}

func NewFetchLessonsHandler(s Service) *FetchLessonsHandler {
	return &FetchLessonsHandler{s: s}
}

func (h *FetchLessonsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	accessToken := req.Header.Get("authorization")
	if accessToken == "" {
		log.Print("missing authorization header")
		res.WriteHeader(http.StatusUnauthorized)

		return
	}

	userID, err := auth.UserID(accessToken)
	if err != nil {
		log.Print(fmt.Errorf("failed to get user ID from access token: %w", err))
		res.WriteHeader(http.StatusUnauthorized)

		return
	}

	lessons, err := h.s.FetchLessons(ctx, userID)
	if err != nil {
		log.Print(fmt.Errorf("failed to fetch lessons: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoded, err := json.Marshal(lessons)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode FetchLessonsHandler HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write FetchLessonsHandler HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
