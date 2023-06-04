package rest

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestCall(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		buf := new(bytes.Buffer)
		_, err := buf.ReadFrom(r.Body)
		assert.NoError(t, err)
		reqBody := buf.String()

		assert.Equal(t, "request body", reqBody)
		assert.Equal(t, "/route?foo=bar", r.URL.String())
		assert.Equal(t, "application/json", r.Header.Get("content-type"))

		w.WriteHeader(http.StatusOK)
		_, err = w.Write([]byte(`response body`))
		assert.NoError(t, err)
	}))
	defer server.Close()

	u, err := url.Parse(server.URL)
	assert.NoError(t, err)

	c := NewClient(server.Client(), u.Host, u.Scheme)

	respBody, err := c.Call("method", "/route", map[string]string{"foo": "bar"}, []byte("request body"))

	assert.NoError(t, err)
	assert.Equal(t, "response body", string(respBody))
}
