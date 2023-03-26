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

func TestFetchNextExercise(t *testing.T) {
	exercise := models.Exercise{
		Id:          1,
		Question:    "question",
		Answer:      "answer",
		BadAnswers:  2,
		GoodAnswers: 3,
	}

	reader := storage.NewReaderMock()
	reader.On("RandomExercise").Return(exercise)

	route := NewFetchNextExercise(reader)

	res := httptest.NewRecorder()
	req := &http.Request{}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var result models.Exercise
	json.Unmarshal(res.Body.Bytes(), &result)

	assert.Equal(t, exercise, result)
}
