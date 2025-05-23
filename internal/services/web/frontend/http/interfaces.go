package http

import "context"

// Caller is an interface for making API calls
type Caller interface {
	Call(ctx context.Context, method, route string, params map[string]string, reqBody []byte, accessToken string) ([]byte, error)
}
