package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/auth"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/validation"
)

type FetchExercisesOfLessonHandler struct {
	s Service
}

func NewFetchExercisesOfLessonHandler(s Service) *FetchExercisesOfLessonHandler {
	return &FetchExercisesOfLessonHandler{s: s}
}

func (h *FetchExercisesOfLessonHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	accessToken := req.Header.Get("authorization")
	if accessToken == "" {
		log.Print("missing authorization header")
		res.WriteHeader(http.StatusUnauthorized)

		return
	}

	userID, err := auth.UserID(accessToken)
	if err != nil {
		log.Print(fmt.Errorf("failed to get user ID from access token: %w", err))
		res.WriteHeader(http.StatusUnauthorized)

		return
	}

	lessonId, err := strconv.Atoi(req.URL.Query().Get("lesson_id"))
	if err != nil {
		log.Print(fmt.Errorf("failed to get a lesson_id: %w", err))

		// validate empty lesson if lesson_id is not present, this is for error messages consistency
		validator := validation.ValidateLessonIdentified(backend.Lesson{})
		if validator.Failed() {
			log.Print(fmt.Errorf("invalid input: %w", validator))

			res.WriteHeader(http.StatusBadRequest)

			encoded, err := json.Marshal(validator.Error())
			if err != nil {
				log.Print(fmt.Errorf("failed to encode FetchExercises HTTP response: %w", err))
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			_, err = res.Write(encoded)
			if err != nil {
				log.Print(fmt.Errorf("failed to write FetchExercises HTTP response: %w", err))
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			return
		}
	}

	oldestExerciseID, err := strconv.Atoi(req.URL.Query().Get("oldest_exercise_id"))
	if err != nil {
		log.Print(fmt.Errorf("failed to get a oldest_exercise_id: %w", err))
	}

	exercises, err := h.s.FetchExercises(ctx, backend.Lesson{Id: lessonId}, oldestExerciseID, userID)
	if err != nil {
		log.Print(fmt.Errorf("failed to fetch exercises: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoded, err := json.Marshal(exercises)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode FetchExercises HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write FetchExercises HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
