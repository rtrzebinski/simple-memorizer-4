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

type UpsertExercise struct {
	w        storage.Writer
	exercise models.Exercise
}

func NewUpsertExercise(w storage.Writer) *UpsertExercise {
	return &UpsertExercise{w: w}
}

func (h *UpsertExercise) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode UpsertExercise HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateUpsertExercise(h.exercise, nil)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode UpsertExercise HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write UpsertExercise HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.w.UpsertExercise(&h.exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to upsert exercise: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
