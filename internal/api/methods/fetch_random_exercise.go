package methods

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-go/internal/storage"
	"log"
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
		log.Fatal(err)
	}

	w.Write(encoded)
}
