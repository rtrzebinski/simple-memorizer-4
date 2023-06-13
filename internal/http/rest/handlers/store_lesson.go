package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/validation"
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

	validator := validation.ValidateStoreLesson(h.lesson, nil)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
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
