package rest

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal"
	"github.com/rtrzebinski/simple-memorizer-4/internal/validation"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestStoreAnswer(t *testing.T) {
	input := models.Answer{
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
	writer.On("StoreAnswer", &input)

	route := NewStoreAnswer(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	writer.AssertExpectations(t)
}

func TestStoreAnswer_invalidInput(t *testing.T) {
	input := models.Answer{}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := internal.NewWriterMock()

	route := NewStoreAnswer(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	assert.Equal(t, http.StatusBadRequest, res.Code)

	var result string

	err = json.Unmarshal(res.Body.Bytes(), &result)
	assert.NoError(t, err)
	assert.Equal(t, validation.ValidateStoreAnswer(models.Answer{}).Error(), result)
}
