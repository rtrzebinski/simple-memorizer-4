package server

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/stretchr/testify/assert"
)

func TestFetchLessonsHandler(t *testing.T) {
	lesson := backend.Lesson{}
	lessons := backend.Lessons{lesson}

	reader := NewReaderMock()
	reader.On("FetchLessons").Return(lessons)

	route := NewFetchLessonsHandler(reader)

	res := httptest.NewRecorder()
	req := &http.Request{}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var result backend.Lessons
	err := json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, lessons, result)
}
