package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"log"
	"net/http"
)

type StoreExercise struct {
	w        storage.Writer
	exercise models.Exercise
}

func NewStoreExercise(w storage.Writer) *StoreExercise {
	return &StoreExercise{w: w}
}

func (h *StoreExercise) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode StoreExercise HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	err = h.w.StoreExercise(h.exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to store exercise: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
