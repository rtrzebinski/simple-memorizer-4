package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/validation"
	"log"
	"net/http"
)

type StoreExercises struct {
	w         storage.Writer
	exercises models.Exercises
}

func NewStoreExercises(w storage.Writer) *StoreExercises {
	return &StoreExercises{w: w}
}

func (h *StoreExercises) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.exercises)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode StoreExercises HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateStoreExercises(h.exercises)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode StoreExercises HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write StoreExercises HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.w.StoreExercises(h.exercises)
	if err != nil {
		log.Print(fmt.Errorf("failed to store exercises: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
