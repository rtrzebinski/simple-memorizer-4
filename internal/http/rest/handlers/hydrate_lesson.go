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

type HydrateLesson struct {
	r storage.Reader
}

func NewHydrateLesson(r storage.Reader) *HydrateLesson {
	return &HydrateLesson{r: r}
}

func (h *HydrateLesson) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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
				log.Print(fmt.Errorf("failed to encode HydrateLesson HTTP response: %w", err))
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			_, err = res.Write(encoded)
			if err != nil {
				log.Print(fmt.Errorf("failed to write HydrateLesson HTTP response: %w", err))
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			return
		}
	}

	lesson := &models.Lesson{Id: lessonId}

	err = h.r.HydrateLesson(lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to hydrate a lesson: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoded, err := json.Marshal(lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode HydrateLesson HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write HydrateLesson HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
