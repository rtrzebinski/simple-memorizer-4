package methods

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-go/internal/storage"
	"log"
	"net/http"
)

type IncrementBadAnswers struct {
	w storage.Writer
}

func NewIncrementBadAnswers(w storage.Writer) *IncrementBadAnswers {
	return &IncrementBadAnswers{w: w}
}

func (h *IncrementBadAnswers) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		ExerciseId int
	}

	var input Input

	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		w.WriteHeader(400)
		log.Println(err.Error())
	}

	h.w.IncrementBadAnswers(input.ExerciseId)
}
