package server

import (
	"context"
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

func TestUpsertExerciseHandler(t *testing.T) {
	ctx := context.Background()

	input := backend.Exercise{
		Question: "question",
		Answer:   "answer",
		Lesson:   &backend.Lesson{Id: 10},
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()
	service.On("UpsertExercise", ctx, &input).Return(nil)

	route := NewUpsertExerciseHandler(service)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	service.AssertExpectations(t)
}

func TestUpsertExerciseHandler_invalidInput(t *testing.T) {
	input := backend.Exercise{}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()

	route := NewUpsertExerciseHandler(service)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateUpsertExercise(input, nil).Error(), result)
}
