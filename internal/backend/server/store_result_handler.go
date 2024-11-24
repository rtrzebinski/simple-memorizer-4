package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/server/validation"
)

type StoreResultHandler struct {
	s Service
}

func NewStoreResultHandler(s Service) *StoreResultHandler {
	return &StoreResultHandler{s: s}
}

func (h *StoreResultHandler) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	var result backend.Result

	// Derive ctx from the request context
	ctx := req.Context()

	err := json.NewDecoder(req.Body).Decode(&result)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode StoreResultHandler HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateStoreResult(result)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode StoreResultHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write StoreResultHandler HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	switch result.Type {
	case backend.Good:
		err = h.s.PublishGoodAnswer(ctx, result.Exercise.Id)
		if err != nil {
			log.Print(fmt.Errorf("failed to publish good answer event: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}
	case backend.Bad:
		err = h.s.PublishBadAnswer(ctx, result.Exercise.Id)
		if err != nil {
			log.Print(fmt.Errorf("failed to publish bad answer event: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}
	}
}
