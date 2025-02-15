package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/server/validation"
)

type StoreExercisesHandler struct {
	s Service
}

func NewStoreExercisesHandler(s Service) *StoreExercisesHandler {
	return &StoreExercisesHandler{s: s}
}

func (h *StoreExercisesHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	var exercises backend.Exercises

	err := json.NewDecoder(req.Body).Decode(&exercises)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode StoreExercisesHandler HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateStoreExercises(exercises)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode StoreExercisesHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write StoreExercisesHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.s.StoreExercises(ctx, exercises)
	if err != nil {
		log.Print(fmt.Errorf("failed to store exercises: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
