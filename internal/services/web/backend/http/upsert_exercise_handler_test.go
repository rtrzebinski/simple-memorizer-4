package http

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/http/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpsertExerciseHandler(t *testing.T) {
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
	service.On("UpsertExercise", mock.Anything, &input, "100").Return(nil)

	v := NewTokenVerifierMock()
	v.On("VerifyAndUserID", mock.Anything, mock.Anything).Return("100", nil)
	route := RequireAuth(v)(NewUpsertExerciseHandler(service))

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}
	req.Header = make(map[string][]string)
	// { "sub": "100" }
	req.Header.Set("authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMDAifQ.bEOa2kaRwC1f7Ow-7WgSltYq-Vz9JUDCo3EPe7KEXd8")

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

	v := NewTokenVerifierMock()
	v.On("VerifyAndUserID", mock.Anything, mock.Anything).Return("100", nil)
	route := RequireAuth(v)(NewUpsertExerciseHandler(service))

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}
	req.Header = make(map[string][]string)
	// { "sub": "100" }
	req.Header.Set("authorization", "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMDAifQ.bEOa2kaRwC1f7Ow-7WgSltYq-Vz9JUDCo3EPe7KEXd8")

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateUpsertExercise(input, nil).Error(), result)
}
