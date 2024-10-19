package server

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/validation"
	"log"
	"net/http"
)

type UpsertLesson struct {
	w      storage.Writer
	lesson models.Lesson
}

func NewUpsertLesson(w storage.Writer) *UpsertLesson {
	return &UpsertLesson{w: w}
}

func (h *UpsertLesson) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode UpsertLesson HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateUpsertLesson(h.lesson, nil)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode UpsertLesson HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write UpsertLesson HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.w.UpsertLesson(&h.lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to upsert lesson: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
