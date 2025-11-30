package http

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestFetchLessonsHandler(t *testing.T) {
	lesson := backend.Lesson{}
	lessons := backend.Lessons{lesson}

	service := NewServiceMock()
	service.On("FetchLessons", mock.Anything, "100").Return(lessons, nil)

	v := NewTokenVerifierMock()
	v.On("VerifyAndUserID", mock.Anything, mock.Anything).Return("100", nil)
	route := RequireAuth(v)(NewFetchLessonsHandler(service))

	res := httptest.NewRecorder()
	req := &http.Request{}
	req.Header = make(map[string][]string)
	// { "sub": "100" }
	req.Header.Set("authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMDAifQ.bEOa2kaRwC1f7Ow-7WgSltYq-Vz9JUDCo3EPe7KEXd8")

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Code)

	var result backend.Lessons
	err := json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)

	assert.Equal(t, lessons, result)
}
