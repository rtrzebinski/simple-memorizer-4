package routes

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"log"
	"net/http"
)

type FetchNextExercise struct {
	r storage.Reader
}

func NewFetchNextExercise(r storage.Reader) *FetchNextExercise {
	return &FetchNextExercise{r: r}
}

func (h *FetchNextExercise) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	exercise, err := h.r.RandomExercise()
	if err != nil {
		log.Print(fmt.Errorf("failed to find a random exercise: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoded, err := json.Marshal(exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode FetchNextExercise HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write FetchNextExercise HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
