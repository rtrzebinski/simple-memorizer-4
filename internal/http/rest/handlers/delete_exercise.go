package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/validation"
	"log"
	"net/http"
)

type DeleteExercise struct {
	w        storage.Writer
	exercise models.Exercise
}

func NewDeleteExercise(w storage.Writer) *DeleteExercise {
	return &DeleteExercise{w: w}
}

func (h *DeleteExercise) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode StoreExercise HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateExerciseIdentified(h.exercise)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode DeleteExercise HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write DeleteExercise HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.w.DeleteExercise(h.exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to store exercise: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
