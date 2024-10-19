package server

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/validation"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteLesson(t *testing.T) {
	input := models.Lesson{
		Id: 123,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := NewWriterMock()
	writer.On("DeleteLesson", input)

	route := NewDeleteLesson(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	writer.AssertExpectations(t)
}

func TestDeleteLesson_invalidInput(t *testing.T) {
	input := models.Lesson{}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := NewWriterMock()

	route := NewDeleteLesson(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateLessonIdentified(input).Error(), result)
}
