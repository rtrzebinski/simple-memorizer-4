package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/auth"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/validation"
)

type HandlerStoreExercises struct {
	s Service
}

func NewHandlerStoreExercises(s Service) *HandlerStoreExercises {
	return &HandlerStoreExercises{s: s}
}

func (h *HandlerStoreExercises) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if req.Method != http.MethodPost {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.UserIDFromContext(req.Context())
	if !ok || userID == "" {
		http.Error(res, "unauthorized", http.StatusUnauthorized)
		return
	}

	var exercises backend.Exercises

	err := json.NewDecoder(req.Body).Decode(&exercises)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode HandlerStoreExercises HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateStoreExercises(exercises)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode HandlerStoreExercises HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write HandlerStoreExercises HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.s.StoreExercises(ctx, userID, exercises)
	if err != nil {
		log.Print(fmt.Errorf("failed to store exercises: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
