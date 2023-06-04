package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"log"
	"net/http"
)

type DeleteLesson struct {
	w      storage.Writer
	lesson models.Lesson
}

func NewDeleteLesson(w storage.Writer) *DeleteLesson {
	return &DeleteLesson{w: w}
}

func (h *DeleteLesson) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode StoreLesson HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	err = h.w.DeleteLesson(h.lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to store lesson: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
