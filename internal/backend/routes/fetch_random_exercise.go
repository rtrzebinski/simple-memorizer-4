package routes

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"log"
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
		log.Print(fmt.Errorf("failed to encode FetchRandomExercise HTTP response: %w", err))
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write FetchRandomExercise HTTP response: %w", err))
	}
}
