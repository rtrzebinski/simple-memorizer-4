package http

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

// Caller is an interface for making API calls
type Caller interface {
	Call(ctx app.Context, method, route string, params map[string]string, reqBody []byte) ([]byte, error)
}
