package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/server/validation"
)

type UpsertLessonHandler struct {
	s Service
}

func NewUpsertLessonHandler(s Service) *UpsertLessonHandler {
	return &UpsertLessonHandler{s: s}
}

func (h *UpsertLessonHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var lesson backend.Lesson

	err := json.NewDecoder(req.Body).Decode(&lesson)
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

	err = h.s.UpsertLesson(&lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to upsert lesson: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
