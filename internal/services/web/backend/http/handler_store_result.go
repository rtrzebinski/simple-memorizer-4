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

type HandlerStoreResult struct {
	s Service
}

func NewHandlerStoreResult(s Service) *HandlerStoreResult {
	return &HandlerStoreResult{s: s}
}

func (h *HandlerStoreResult) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if req.Method != http.MethodPost {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := auth.UserIDFromContext(ctx)
	if !ok || userID == "" {
		http.Error(res, "unauthorized", http.StatusUnauthorized)
		return
	}

	var result backend.Result

	err := json.NewDecoder(req.Body).Decode(&result)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode HandlerStoreResult HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateStoreResult(result)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode HandlerStoreResult HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write HandlerStoreResult HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	switch result.Type {
	case backend.Good:
		err = h.s.PublishGoodAnswer(ctx, userID, result.Exercise.Id)
		if err != nil {
			log.Print(fmt.Errorf("failed to publish good answer event: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}
	case backend.Bad:
		err = h.s.PublishBadAnswer(ctx, userID, result.Exercise.Id)
		if err != nil {
			log.Print(fmt.Errorf("failed to publish bad answer event: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}
	}
}
