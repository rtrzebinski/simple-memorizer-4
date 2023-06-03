package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"log"
	"net/http"
)

type IncrementGoodAnswers struct {
	w        storage.Writer
	exercise models.Exercise
}

func NewIncrementGoodAnswers(w storage.Writer) *IncrementGoodAnswers {
	return &IncrementGoodAnswers{w: w}
}

func (h *IncrementGoodAnswers) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode IncrementGoodAnswers HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	err = h.w.IncrementGoodAnswers(h.exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to increment bad answers: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
