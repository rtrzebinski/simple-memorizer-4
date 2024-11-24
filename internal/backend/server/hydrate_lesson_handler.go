package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/server/validation"
)

type HydrateLessonHandler struct {
	s Service
}

func NewHydrateLessonHandler(s Service) *HydrateLessonHandler {
	return &HydrateLessonHandler{s: s}
}

func (h *HydrateLessonHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	lessonId, err := strconv.Atoi(req.URL.Query().Get("lesson_id"))
	if err != nil {
		log.Print(fmt.Errorf("failed to get a lesson_id: %w", err))

		// validate empty lesson if lesson_id is not present, this is for error messages consistency
		validator := validation.ValidateLessonIdentified(backend.Lesson{})
		if validator.Failed() {
			log.Print(fmt.Errorf("invalid input: %w", validator))

			res.WriteHeader(http.StatusBadRequest)

			encoded, err := json.Marshal(validator.Error())
			if err != nil {
				log.Print(fmt.Errorf("failed to encode HydrateLessonHandler HTTP response: %w", err))
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			_, err = res.Write(encoded)
			if err != nil {
				log.Print(fmt.Errorf("failed to write HydrateLessonHandler HTTP response: %w", err))
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			return
		}
	}

	lesson := &backend.Lesson{Id: lessonId}

	err = h.s.HydrateLesson(ctx, lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to hydrate a lesson: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoded, err := json.Marshal(lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode HydrateLessonHandler HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write HydrateLessonHandler HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
