package routes

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchAllExercises(t *testing.T) {
	exercise := models.Exercise{
		Id:          1,
		Question:    "question",
		Answer:      "answer",
		BadAnswers:  2,
		GoodAnswers: 3,
	}
	exercises := models.Exercises{exercise}

	reader := storage.NewReaderMock()
	reader.On("AllExercises").Return(exercises)

	route := NewFetchAllExercises(reader)

	res := httptest.NewRecorder()
	req := &http.Request{}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var result models.Exercises
	json.Unmarshal(res.Body.Bytes(), &result)

	assert.Equal(t, exercises, result)
}
