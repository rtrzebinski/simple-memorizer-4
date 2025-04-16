package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/validation"
	"github.com/stretchr/testify/assert"
)

func TestUpsertLessonHandler(t *testing.T) {
	ctx := context.Background()

	input := backend.Lesson{
		Name:        "name",
		Description: "description",
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()
	service.On("UpsertLesson", ctx, &input, "100").Return(nil)

	route := NewUpsertLessonHandler(service)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}
	req.Header = make(map[string][]string)
	// { "sub": "100" }
	req.Header.Set("authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMDAifQ.bEOa2kaRwC1f7Ow-7WgSltYq-Vz9JUDCo3EPe7KEXd8")

	route.ServeHTTP(res, req)

	service.AssertExpectations(t)
}

func TestUpsertLessonHandler_invalidInput(t *testing.T) {
	input := backend.Lesson{}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()

	route := NewUpsertLessonHandler(service)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}
	req.Header = make(map[string][]string)
	// { "sub": "100" }
	req.Header.Set("authorization", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMDAifQ.bEOa2kaRwC1f7Ow-7WgSltYq-Vz9JUDCo3EPe7KEXd8")

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateUpsertLesson(input, nil).Error(), result)
}
