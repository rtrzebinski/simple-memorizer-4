package routes

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"log"
	"net/http"
)

type StoreLesson struct {
	w      storage.Writer
	lesson models.Lesson
}

func NewStoreLesson(w storage.Writer) *StoreLesson {
	return &StoreLesson{w: w}
}

func (h *StoreLesson) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode StoreLesson HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	err = h.w.StoreLesson(h.lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to store lesson: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}