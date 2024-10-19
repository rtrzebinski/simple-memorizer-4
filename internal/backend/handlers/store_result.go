package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/validation"
	"log"
	"net/http"
)

type StoreResult struct {
	w      storage.Writer
	result models.Result
}

func NewStoreResult(w storage.Writer) *StoreResult {
	return &StoreResult{w: w}
}

func (h *StoreResult) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.result)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode StoreResult HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateStoreResult(h.result)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode StoreResult HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write StoreResult HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.w.StoreResult(&h.result)
	if err != nil {
		log.Print(fmt.Errorf("failed to store result: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
