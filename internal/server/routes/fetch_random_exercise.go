package routes

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/server/storage"
	"net/http"
)

type FetchRandomExercise struct {
	r storage.Reader
}

func NewFetchRandomExercise(r storage.Reader) *FetchRandomExercise {
	return &FetchRandomExercise{r: r}
}

func (h *FetchRandomExercise) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	exercise := h.r.RandomExercise()

	encoded, err := json.Marshal(exercise)
	if err != nil {
		panic(err)
	}

	_, err = w.Write(encoded)
	if err != nil {
		panic(err)
	}
}
