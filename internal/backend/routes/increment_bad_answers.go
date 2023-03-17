package routes

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
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
		panic(err)
	}

	h.w.IncrementBadAnswers(h.req.ExerciseId)
}
