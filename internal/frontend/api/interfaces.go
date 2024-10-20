package api

// Caller is an interface for making API calls
type Caller interface {
	Call(method, route string, params map[string]string, reqBody []byte) ([]byte, error)
}
