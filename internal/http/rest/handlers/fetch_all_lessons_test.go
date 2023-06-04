package rest

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestFetchAllLessons(t *testing.T) {
	lesson := models.Lesson{
		Id:   1,
		Name: "name",
	}
	lessons := models.Lessons{lesson}

	reader := storage.NewReaderMock()
	reader.On("FetchAllLessons").Return(lessons)

	route := NewFetchAllLessons(reader)

	res := httptest.NewRecorder()
	req := &http.Request{}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var result models.Lessons
	json.Unmarshal(res.Body.Bytes(), &result)

	assert.Equal(t, lessons, result)
}