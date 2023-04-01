package routes

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIncrementGoodAnswers(t *testing.T) {
	input := models.Exercise{Id: 456}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := storage.NewWriterMock()
	writer.On("IncrementGoodAnswers", 456)

	route := NewIncrementGoodAnswers(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	writer.AssertExpectations(t)
}
