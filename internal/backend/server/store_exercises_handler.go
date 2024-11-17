package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/server/validation"
)

type StoreExercisesHandler struct {
	s         Service
	exercises backend.Exercises
}

func NewStoreExercisesHandler(s Service) *StoreExercisesHandler {
	return &StoreExercisesHandler{s: s}
}

func (h *StoreExercisesHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.exercises)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode StoreExercisesHandler HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateStoreExercises(h.exercises)
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

	err = h.s.StoreExercises(h.exercises)
	if err != nil {
		log.Print(fmt.Errorf("failed to store exercises: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
