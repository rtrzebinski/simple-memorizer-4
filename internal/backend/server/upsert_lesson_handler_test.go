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

func TestUpsertLessonHandler(t *testing.T) {
	input := backend.Lesson{
		Name:        "name",
		Description: "description",
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := NewWriterMock()
	writer.On("UpsertLesson", &input)

	route := NewUpsertLessonHandler(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	writer.AssertExpectations(t)
}

func TestUpsertLessonHandler_invalidInput(t *testing.T) {
	input := backend.Lesson{}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := NewWriterMock()

	route := NewUpsertLessonHandler(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateUpsertLesson(input, nil).Error(), result)
}
