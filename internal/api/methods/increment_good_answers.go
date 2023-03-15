package methods

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-go/internal/storage"
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

func (h *IncrementGoodAnswers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := json.NewDecoder(r.Body).Decode(&h.req)
	if err != nil {
		panic(err)
	}

	h.w.IncrementGoodAnswers(h.req.ExerciseId)
}
