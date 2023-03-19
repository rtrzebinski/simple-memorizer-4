package routes

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"log"
	"net/http"
)

type IncrementBadAnswers struct {
	w   storage.Writer
	req IncrementBadAnswersReq
}

type IncrementBadAnswersReq struct {
	ExerciseId int
}

func NewIncrementBadAnswers(w storage.Writer) *IncrementBadAnswers {
	return &IncrementBadAnswers{w: w}
}

func (h *IncrementBadAnswers) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.req)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode IncrementBadAnswers HTTP request: %w", err))
	}

	h.w.IncrementBadAnswers(h.req.ExerciseId)
}
