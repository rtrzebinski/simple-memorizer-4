package api

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type HttpCaller struct {
	http   *http.Client
	host   string
	scheme string
}

func NewHttpCaller(http *http.Client, host string, scheme string) *HttpCaller {
	return &HttpCaller{http: http, host: host, scheme: scheme}
}

func (c *HttpCaller) Call(method, route string, params map[string]string, reqBody []byte) ([]byte, error) {
	// parse url
	u, err := url.Parse(c.scheme + "://" + c.host + route)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTTP request URL: %w", err)
	}

	// encode query params
	if params != nil {
		p := u.Query()
		for k, v := range params {
			p.Add(k, v)
		}
		u.RawQuery = p.Encode()
	}

	// create request body buffer
	buffer := bytes.NewBuffer(reqBody)

	// create a request
	req, err := http.NewRequest(method, u.String(), buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// add content-type header
	req.Header.Add("content-type", "application/json")

	// make a request
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call HTTP endpoint: %w", err)
	}

	// defer body closing
	defer resp.Body.Close()

	// read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	// check if status is OK
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned with the status code '%d': %w",
			resp.StatusCode, errors.New(string(respBody)))
	}

	return respBody, nil
}
