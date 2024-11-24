package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/server/validation"
)

type UpsertExerciseHandler struct {
	s Service
}

func NewUpsertExerciseHandler(s Service) *UpsertExerciseHandler {
	return &UpsertExerciseHandler{s: s}
}

func (h *UpsertExerciseHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var exercise backend.Exercise

	err := json.NewDecoder(req.Body).Decode(&exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode UpsertExerciseHandler HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateUpsertExercise(exercise, nil)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode UpsertExerciseHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write UpsertExerciseHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.s.UpsertExercise(ctx, &exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to upsert exercise: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
