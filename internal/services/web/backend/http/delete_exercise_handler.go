package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/auth"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/validation"
)

type DeleteExerciseHandler struct {
	s Service
}

func NewDeleteExerciseHandler(s Service) *DeleteExerciseHandler {
	return &DeleteExerciseHandler{s: s}
}

func (h *DeleteExerciseHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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

	var exercise backend.Exercise

	err = json.NewDecoder(req.Body).Decode(&exercise)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode DeleteExerciseHandler HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateExerciseIdentified(exercise)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode DeleteExerciseHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write DeleteExerciseHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.s.DeleteExercise(ctx, exercise, userID)
	if err != nil {
		log.Print(fmt.Errorf("failed to delete exercise: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
