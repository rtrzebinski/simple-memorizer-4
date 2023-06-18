package rest

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/validation"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestFetchResultsOfExercise(t *testing.T) {
	answer := models.Result{
		Id:        10,
		Type:      models.Good,
		CreatedAt: time.Time{},
	}
	answers := models.Results{answer}

	exerciseId := 10

	reader := internal.NewReaderMock()
	reader.On("FetchResultsOfExercise", models.Exercise{Id: exerciseId}).Return(answers)

	route := NewFetchResultsOfExercise(reader)

	u, _ := url.Parse("/")
	params := u.Query()
	params.Add("exercise_id", strconv.Itoa(exerciseId))
	u.RawQuery = params.Encode()

	req := &http.Request{}
	req.URL = u

	res := httptest.NewRecorder()

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var result models.Results
	json.Unmarshal(res.Body.Bytes(), &result)

	assert.Equal(t, answers, result)
}

func TestFetchAnswersOfExercise_invalidInput(t *testing.T) {
	reader := internal.NewReaderMock()

	route := NewFetchResultsOfExercise(reader)

	u, _ := url.Parse("/")

	req := &http.Request{}
	req.URL = u

	res := httptest.NewRecorder()

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err := json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateExerciseIdentified(models.Exercise{}).Error(), result)
}
