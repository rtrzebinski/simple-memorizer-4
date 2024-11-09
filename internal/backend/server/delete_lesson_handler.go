package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/server/validation"
)

type DeleteLessonHandler struct {
	w      Writer
	lesson backend.Lesson
}

func NewDeleteLessonHandler(w Writer) *DeleteLessonHandler {
	return &DeleteLessonHandler{w: w}
}

func (h *DeleteLessonHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode DeleteLessonHandler HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateLessonIdentified(h.lesson)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode DeleteLessonHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write DeleteLessonHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.w.DeleteLesson(h.lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to delete lesson: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
