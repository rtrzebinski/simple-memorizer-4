package http

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	bauth "github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/auth"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/validation"
)

type HandlerUpsertLesson struct {
	s Service
}

func NewHandlerUpsertLesson(s Service) *HandlerUpsertLesson {
	return &HandlerUpsertLesson{s: s}
}

func (h *HandlerUpsertLesson) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	if req.Method != http.MethodPost {
		http.Error(res, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := bauth.UserIDFromContext(ctx)
	if !ok || userID == "" {
		http.Error(res, "unauthorized", http.StatusUnauthorized)
		return
	}

	var lesson backend.Lesson

	err := json.NewDecoder(req.Body).Decode(&lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode HandlerUpsertLesson HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateUpsertLesson(lesson, nil)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode HandlerUpsertLesson HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write HandlerUpsertLesson HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.s.UpsertLesson(ctx, userID, &lesson)
	if err != nil {
		log.Print(fmt.Errorf("failed to upsert lesson: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
