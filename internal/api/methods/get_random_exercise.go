package methods

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-go/internal/storage"
	"log"
	"net/http"
)

type GetRandomExercise struct {
	r storage.Reader
}

func NewExercisesHandler(r storage.Reader) *GetRandomExercise {
	return &GetRandomExercise{
		r: r,
	}
}

func (h *GetRandomExercise) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	exercise := h.r.RandomExercise()

	encoded, err := json.Marshal(exercise)
	if err != nil {
		log.Fatal(err)
	}

	w.Write(encoded)
}
