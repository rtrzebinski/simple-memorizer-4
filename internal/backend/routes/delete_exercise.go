package routes

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
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

	err = h.w.DeleteExercise(h.exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to store exercise: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
