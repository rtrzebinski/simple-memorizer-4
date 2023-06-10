package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/validators"
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

	err = validators.ValidateStoreLesson(h.lesson, nil)
	if err != nil {
		log.Print(fmt.Errorf("invalid input: %w", err))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(err.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode StoreLesson HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write StoreLesson HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.w.StoreLesson(&h.lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to store lesson: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
