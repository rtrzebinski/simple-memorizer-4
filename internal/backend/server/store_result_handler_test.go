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

	publisher := NewPublisherMock()
	publisher.On("PublishGoodAnswer", mock.AnythingOfType("context.backgroundCtx"), 10)

	route := NewStoreResultHandler(publisher)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	publisher.AssertExpectations(t)
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

	publisher := NewPublisherMock()
	publisher.On("PublishBadAnswer", mock.AnythingOfType("context.backgroundCtx"), 10)

	route := NewStoreResultHandler(publisher)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	publisher.AssertExpectations(t)
}

func TestStoreResultHandler_invalidInput(t *testing.T) {
	input := backend.Result{}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	publisher := NewPublisherMock()

	route := NewStoreResultHandler(publisher)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateStoreResult(input).Error(), result)
}
