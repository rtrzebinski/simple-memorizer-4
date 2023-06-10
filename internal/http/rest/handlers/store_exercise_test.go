package rest

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStoreExercise(t *testing.T) {
	input := models.Exercise{
		Question: "question",
		Answer:   "answer",
		Lesson:   &models.Lesson{Id: 10},
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := storage.NewWriterMock()
	writer.On("StoreExercise", &input)

	route := NewStoreExercise(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	writer.AssertExpectations(t)
}

func TestStoreExercise_invalidInput(t *testing.T) {
	input := models.Exercise{}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := storage.NewWriterMock()

	route := NewStoreExercise(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)
}
