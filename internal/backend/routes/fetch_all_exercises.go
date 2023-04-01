package routes

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"log"
	"net/http"
)

type FetchAllExercises struct {
	r storage.Reader
}

func NewFetchAllExercises(r storage.Reader) *FetchAllExercises {
	return &FetchAllExercises{r: r}
}

func (h *FetchAllExercises) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	exercises, err := h.r.AllExercises()
	if err != nil {
		log.Print(fmt.Errorf("failed to fetch all exercises: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoded, err := json.Marshal(exercises)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode FetchAllExercises HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write FetchAllExercises HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
