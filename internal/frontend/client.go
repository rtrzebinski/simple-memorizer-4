package frontend

type Client interface {
	Call(method, route string, params map[string]string, reqBody []byte) ([]byte, error)
}
