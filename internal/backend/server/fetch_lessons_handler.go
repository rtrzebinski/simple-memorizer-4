package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type FetchLessonsHandler struct {
	s Service
}

func NewFetchLessonsHandler(s Service) *FetchLessonsHandler {
	return &FetchLessonsHandler{s: s}
}

func (h *FetchLessonsHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	lessons, err := h.s.FetchLessons()
	if err != nil {
		log.Print(fmt.Errorf("failed to fetch lessons: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoded, err := json.Marshal(lessons)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode FetchLessonsHandler HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write FetchLessonsHandler HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
