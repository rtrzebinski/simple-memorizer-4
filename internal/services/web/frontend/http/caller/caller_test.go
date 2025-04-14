package caller

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCaller_Call(t *testing.T) {
	ctx := context.Background()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(r.Body)
		assert.NoError(t, err)
		reqBody := buf.String()

		assert.Equal(t, "request body", reqBody)
		assert.Equal(t, "/route?foo=bar", r.URL.String())
		assert.Equal(t, "application/json", r.Header.Get("content-type"))
		assert.Equal(t, "authToken", r.Header.Get("authorization"))

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(`response body`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	u, err := url.Parse(server.URL)
	assert.NoError(t, err)

	c := NewCaller(server.Client(), u.Host, u.Scheme)

	respBody, err := c.Call(ctx, "method", "/route", map[string]string{"foo": "bar"}, []byte("request body"), "authToken")

	assert.NoError(t, err)
	assert.Equal(t, "response body", string(respBody))
}
