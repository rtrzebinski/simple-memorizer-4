package server

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"log"
	"net/http"
)

type FetchLessons struct {
	r storage.Reader
}

func NewFetchLessons(r storage.Reader) *FetchLessons {
	return &FetchLessons{r: r}
}

func (h *FetchLessons) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	lessons, err := h.r.FetchLessons()
	if err != nil {
		log.Print(fmt.Errorf("failed to fetch lessons: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoded, err := json.Marshal(lessons)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode FetchLessons HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write FetchLessons HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
