package rest

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/validation"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStoreResult(t *testing.T) {
	input := models.Result{
		Exercise: &models.Exercise{
			Id: 10,
		},
		Type: models.Good,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := internal.NewWriterMock()
	writer.On("StoreResult", &input)

	route := NewStoreResult(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	writer.AssertExpectations(t)
}

func TestStoreResult_invalidInput(t *testing.T) {
	input := models.Result{}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := internal.NewWriterMock()

	route := NewStoreResult(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateStoreResult(models.Result{}).Error(), result)
}
