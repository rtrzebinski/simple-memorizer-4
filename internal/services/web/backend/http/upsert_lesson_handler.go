package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/auth"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/validation"
)

type UpsertLessonHandler struct {
	s Service
}

func NewUpsertLessonHandler(s Service) *UpsertLessonHandler {
	return &UpsertLessonHandler{s: s}
}

func (h *UpsertLessonHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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

	var lesson backend.Lesson

	err = json.NewDecoder(req.Body).Decode(&lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode UpsertLessonHandler HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateUpsertLesson(lesson, nil)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode UpsertLessonHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write UpsertLessonHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.s.UpsertLesson(ctx, &lesson, userID)
	if err != nil {
		log.Print(fmt.Errorf("failed to upsert lesson: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
