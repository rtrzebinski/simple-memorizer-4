package rest

import (
	"encoding/json"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/validation"
	"log"
	"net/http"
	"strconv"
)

type FetchResultsOfExercise struct {
	r internal.Reader
}

func NewFetchResultsOfExercise(r internal.Reader) *FetchResultsOfExercise {
	return &FetchResultsOfExercise{r: r}
}

func (h *FetchResultsOfExercise) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	exerciseId, err := strconv.Atoi(req.URL.Query().Get("exercise_id"))
	if err != nil {
		log.Print(fmt.Errorf("failed to get a exercise_id: %w", err))

		// validate empty exercise if exercise_id is not present, this is for error messages consistency
		validator := validation.ValidateExerciseIdentified(models.Exercise{})
		if validator.Failed() {
			log.Print(fmt.Errorf("invalid input: %w", validator))

			res.WriteHeader(http.StatusBadRequest)

			encoded, err := json.Marshal(validator.Error())
			if err != nil {
				log.Print(fmt.Errorf("failed to encode FetchResultsOfExercise HTTP response: %w", err))
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			_, err = res.Write(encoded)
			if err != nil {
				log.Print(fmt.Errorf("failed to write FetchResultsOfExercise HTTP response: %w", err))
				res.WriteHeader(http.StatusInternalServerError)

				return
			}

			return
		}
	}

	results, err := h.r.FetchResultsOfExercise(models.Exercise{Id: exerciseId})
	if err != nil {
		log.Print(fmt.Errorf("failed to fetch results of exercise: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	encoded, err := json.Marshal(results)
	if err != nil {
		log.Print(fmt.Errorf("failed to encode FetchResultsOfExercise HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}

	_, err = res.Write(encoded)
	if err != nil {
		log.Print(fmt.Errorf("failed to write FetchResultsOfExercise HTTP response: %w", err))
		res.WriteHeader(http.StatusInternalServerError)

		return
	}
}
