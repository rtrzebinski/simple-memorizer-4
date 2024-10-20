package server

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/server/validation"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUpsertExerciseHandler(t *testing.T) {
	input := models.Exercise{
		Question: "question",
		Answer:   "answer",
		Lesson:   &models.Lesson{Id: 10},
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := NewWriterMock()
	writer.On("UpsertExercise", &input)

	route := NewUpsertExerciseHandler(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	writer.AssertExpectations(t)
}

func TestUpsertExerciseHandler_invalidInput(t *testing.T) {
	input := models.Exercise{}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := NewWriterMock()

	route := NewUpsertExerciseHandler(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateUpsertExercise(input, nil).Error(), result)
}
