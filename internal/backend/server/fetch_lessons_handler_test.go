package server

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchLessonsHandler(t *testing.T) {
	lesson := models.Lesson{}
	lessons := models.Lessons{lesson}

	reader := NewReaderMock()
	reader.On("FetchLessons").Return(lessons)

	route := NewFetchLessonsHandler(reader)

	res := httptest.NewRecorder()
	req := &http.Request{}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var result models.Lessons
	json.Unmarshal(res.Body.Bytes(), &result)

	assert.Equal(t, lessons, result)
}
