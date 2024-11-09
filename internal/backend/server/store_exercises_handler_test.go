package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/server/validation"
	"github.com/stretchr/testify/assert"
)

func TestStoreExercises(t *testing.T) {
	input := backend.Exercises{
		backend.Exercise{
			Question: "question",
			Answer:   "answer",
		},
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := NewWriterMock()
	writer.On("StoreExercises", input)

	route := NewStoreExercisesHandler(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	writer.AssertExpectations(t)
}

func TestStoreExercisesHandler_invalidInput(t *testing.T) {
	input := backend.Exercises{
		backend.Exercise{},
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := NewWriterMock()

	route := NewStoreExercisesHandler(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateStoreExercises(input).Error(), result)
}
