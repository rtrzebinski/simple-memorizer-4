package routes

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"net/http"
)

type FetchRandomExercise struct {
	r storage.Reader
}

func NewFetchRandomExercise(r storage.Reader) *FetchRandomExercise {
	return &FetchRandomExercise{r: r}
}

func (h *FetchRandomExercise) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	exercise := h.r.RandomExercise()

	encoded, err := json.Marshal(exercise)
	if err != nil {
		panic(err)
	}

	_, err = res.Write(encoded)
	if err != nil {
		panic(err)
	}
}
