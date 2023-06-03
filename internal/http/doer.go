package http

import "net/http"

type Doer interface {
	Do(req *http.Request) (*http.Response, error)
}
