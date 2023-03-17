package routes

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/server/storage"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestIncrementBadAnswers(t *testing.T) {
	input := IncrementBadAnswersReq{
		ExerciseId: 123,
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	writer := storage.NewWriterMock()
	writer.On("IncrementBadAnswers", 123)

	route := NewIncrementBadAnswers(writer)

	res := httptest.NewRecorder()
	req := &http.Request{Body: io.NopCloser(strings.NewReader(string(body)))}

	route.ServeHTTP(res, req)

	writer.AssertExpectations(t)
}