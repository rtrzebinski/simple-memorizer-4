package server

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/server/validation"
	"log"
	"net/http"
)

type DeleteExerciseHandler struct {
	s Service
}

func NewDeleteExerciseHandler(s Service) *DeleteExerciseHandler {
	return &DeleteExerciseHandler{s: s}
}

func (h *DeleteExerciseHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var exercise backend.Exercise

	err := json.NewDecoder(req.Body).Decode(&exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode DeleteExerciseHandler HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateExerciseIdentified(exercise)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode DeleteExerciseHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write DeleteExerciseHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.s.DeleteExercise(ctx, exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to delete exercise: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
