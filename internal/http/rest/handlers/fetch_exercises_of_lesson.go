package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/validators"
	"log"
	"net/http"
	"strconv"
)

type FetchExercisesOfLesson struct {
	r storage.Reader
}

func NewFetchExercisesOfLesson(r storage.Reader) *FetchExercisesOfLesson {
	return &FetchExercisesOfLesson{r: r}
}

func (h *FetchExercisesOfLesson) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	lessonId, err := strconv.Atoi(req.URL.Query().Get("lesson_id"))
	if err != nil {
		log.Print(fmt.Errorf("failed to get a lesson_id: %w", err))

		// validate empty lesson if lesson_id is not present, this is for error messages consistency
		err = validators.ValidateLessonIdentified(models.Lesson{})
		if err != nil {
			log.Print(fmt.Errorf("invalid input: %w", err))

			res.WriteHeader(http.StatusBadRequest)

			encoded, err := json.Marshal(err.Error())
			if err != nil {
				log.Print(fmt.Errorf("failed to encode FetchExercisesOfLesson HTTP response: %w", err))
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			_, err = res.Write(encoded)
			if err != nil {
				log.Print(fmt.Errorf("failed to write FetchExercisesOfLesson HTTP response: %w", err))
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			return
		}
	}

	exercises, err := h.r.FetchExercisesOfLesson(models.Lesson{Id: lessonId})
	if err != nil {
		log.Print(fmt.Errorf("failed to fetch all exercises: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoded, err := json.Marshal(exercises)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode FetchExercisesOfLesson HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write FetchExercisesOfLesson HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
