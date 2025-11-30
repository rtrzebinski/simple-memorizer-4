package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/auth"
)

type FetchLessonsHandler struct {
	s Service
}

func NewFetchLessonsHandler(s Service) *FetchLessonsHandler {
	return &FetchLessonsHandler{s: s}
}

func (h *FetchLessonsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok || userID == "" {
		http.Error(res, "unauthorized", http.StatusUnauthorized)
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
