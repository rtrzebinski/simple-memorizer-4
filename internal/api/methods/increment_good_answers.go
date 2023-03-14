package methods

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-go/internal/storage"
	"log"
	"net/http"
)

type IncrementGoodAnswers struct {
	w storage.Writer
}

func NewIncrementGoodAnswers(w storage.Writer) *IncrementGoodAnswers {
	return &IncrementGoodAnswers{w: w}
}

func (h *IncrementGoodAnswers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		ExerciseId int
	}

	var input Input

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(400)
		log.Println(err.Error())
	}

	h.w.IncrementGoodAnswers(input.ExerciseId)
}
