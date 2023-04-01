package routes

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"log"
	"net/http"
)

type FetchExercises struct {
	r storage.Reader
}

func NewFetchExercises(r storage.Reader) *FetchExercises {
	return &FetchExercises{r: r}
}

func (h *FetchExercises) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	exercises, err := h.r.Exercises()
	if err != nil {
		log.Print(fmt.Errorf("failed to fetch exercises: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoded, err := json.Marshal(exercises)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode FetchExercises HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write FetchExercises HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
