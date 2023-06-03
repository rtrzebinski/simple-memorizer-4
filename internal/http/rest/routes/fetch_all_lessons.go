package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"log"
	"net/http"
)

type FetchAllLessons struct {
	r storage.Reader
}

func NewFetchAllLessons(r storage.Reader) *FetchAllLessons {
	return &FetchAllLessons{r: r}
}

func (h *FetchAllLessons) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	lessons, err := h.r.AllLessons()
	if err != nil {
		log.Print(fmt.Errorf("failed to fetch all lessons: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoded, err := json.Marshal(lessons)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode FetchAllLessons HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write FetchAllLessons HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
