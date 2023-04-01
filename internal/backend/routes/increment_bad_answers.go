package routes

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"log"
	"net/http"
)

type IncrementBadAnswers struct {
	w        storage.Writer
	exercise models.Exercise
}

func NewIncrementBadAnswers(w storage.Writer) *IncrementBadAnswers {
	return &IncrementBadAnswers{w: w}
}

func (h *IncrementBadAnswers) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode IncrementBadAnswers HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	err = h.w.IncrementBadAnswers(h.exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to increment bad answers: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
