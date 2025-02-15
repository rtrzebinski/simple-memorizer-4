package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/stretchr/testify/assert"
)

func TestFetchLessonsHandler(t *testing.T) {
	ctx := context.Background()

	lesson := backend.Lesson{}
	lessons := backend.Lessons{lesson}

	service := NewServiceMock()
	service.On("FetchLessons", ctx).Return(lessons, nil)

	route := NewFetchLessonsHandler(service)

	res := httptest.NewRecorder()
	req := &http.Request{}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var result backend.Lessons
	err := json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, lessons, result)
}
