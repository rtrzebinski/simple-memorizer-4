package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/validation"
	"log"
	"net/http"
)

type StoreAnswer struct {
	w      internal.Writer
	answer models.Answer
}

func NewStoreAnswer(w internal.Writer) *StoreAnswer {
	return &StoreAnswer{w: w}
}

func (h *StoreAnswer) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	err := json.NewDecoder(req.Body).Decode(&h.answer)
	if err != nil {
		log.Print(fmt.Errorf("failed to decode StoreAnswer HTTP request: %w", err))
		res.WriteHeader(http.StatusBadRequest)

		return
	}

	validator := validation.ValidateStoreAnswer(h.answer)
	if validator.Failed() {
		log.Print(fmt.Errorf("invalid input: %w", validator))

		res.WriteHeader(http.StatusBadRequest)

		encoded, err := json.Marshal(validator.Error())
		if err != nil {
			log.Print(fmt.Errorf("failed to encode StoreAnswer HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		_, err = res.Write(encoded)
		if err != nil {
			log.Print(fmt.Errorf("failed to write StoreAnswer HTTP response: %w", err))
			res.WriteHeader(http.StatusInternalServerError)

			return
		}

		return
	}

	err = h.w.StoreAnswer(&h.answer)
	if err != nil {
		log.Print(fmt.Errorf("failed to store answer: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
