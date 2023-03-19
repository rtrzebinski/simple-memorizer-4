package routes

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"log"
	"net/http"
)

type IncrementGoodAnswers struct {
	w   storage.Writer
	req IncrementGoodAnswersReq
}

type IncrementGoodAnswersReq struct {
	ExerciseId int
}

func NewIncrementGoodAnswers(w storage.Writer) *IncrementGoodAnswers {
	return &IncrementGoodAnswers{w: w}
}

func (h *IncrementGoodAnswers) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.req)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode IncrementGoodAnswers HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	err = h.w.IncrementGoodAnswers(h.req.ExerciseId)
	if err != nil {
		log.Print(fmt.Errorf("failed to increment bad answers: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
