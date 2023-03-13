package handlers

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-go/internal/storage"
	"log"
	"net/http"
)

type ExercisesHandler struct {
	r storage.Reader
}

func NewExercisesHandler(r storage.Reader) *ExercisesHandler {
	return &ExercisesHandler{
		r: r,
	}
}

func (h *ExercisesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	exercise := h.r.RandomExercise()

	encoded, err := json.Marshal(exercise)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(encoded)
}
