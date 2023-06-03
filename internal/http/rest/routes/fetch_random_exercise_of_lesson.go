package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"log"
	"net/http"
	"strconv"
)

type FetchRandomExerciseOfLesson struct {
	r storage.Reader
}

func NewFetchRandomExerciseOfLesson(r storage.Reader) *FetchRandomExerciseOfLesson {
	return &FetchRandomExerciseOfLesson{r: r}
}

func (h *FetchRandomExerciseOfLesson) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	lessonId, err := strconv.Atoi(req.URL.Query().Get("lesson_id"))
	if err != nil {
		log.Print(fmt.Errorf("failed to get a lesson_id: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	exercise, err := h.r.FetchRandomExerciseOfLesson(models.Lesson{Id: lessonId})
	if err != nil {
		log.Print(fmt.Errorf("failed to find a random exercise: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoded, err := json.Marshal(exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode FetchRandomExerciseOfLesson HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write FetchRandomExerciseOfLesson HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}