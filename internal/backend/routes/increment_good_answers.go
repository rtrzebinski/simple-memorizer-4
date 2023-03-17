package routes

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
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
		panic(err)
	}

	h.w.IncrementGoodAnswers(h.req.ExerciseId)
}
