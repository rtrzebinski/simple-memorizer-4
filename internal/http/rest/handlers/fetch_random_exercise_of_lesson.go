package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/validation"
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

		// validate empty lesson if lesson_id is not present, this is for error messages consistency
		validator := validation.ValidateLessonIdentified(models.Lesson{})
		if validator.Failed() {
			log.Print(fmt.Errorf("invalid input: %w", validator))

			res.WriteHeader(http.StatusBadRequest)

			encoded, err := json.Marshal(validator.Error())
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

			return
		}
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
