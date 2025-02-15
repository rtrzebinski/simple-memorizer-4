package server

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend/server/validation"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestStoreResultHandler_goodAnswer(t *testing.T) {
	input := backend.Result{
		Exercise: &backend.Exercise{
			Id: 10,
		},
		Type: backend.Good,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()
	service.On("PublishGoodAnswer", mock.AnythingOfType("context.backgroundCtx"), 10).Return(nil)

	route := NewStoreResultHandler(service)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	service.AssertExpectations(t)
}

func TestStoreResultHandler_badAnswer(t *testing.T) {
	input := backend.Result{
		Exercise: &backend.Exercise{
			Id: 10,
		},
		Type: backend.Bad,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()
	service.On("PublishBadAnswer", mock.AnythingOfType("context.backgroundCtx"), 10).Return(nil)

	route := NewStoreResultHandler(service)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	service.AssertExpectations(t)
}

func TestStoreResultHandler_invalidInput(t *testing.T) {
	input := backend.Result{}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	service := NewServiceMock()

	route := NewStoreResultHandler(service)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateStoreResult(input).Error(), result)
}
