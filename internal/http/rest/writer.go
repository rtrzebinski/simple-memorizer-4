package rest

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	myhttp "github.com/rtrzebinski/simple-memorizer-4/internal/http"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"io"
	"net/http"
	"net/url"
)

type Writer struct {
	http   myhttp.Doer
	host   string
	scheme string
}

func NewWriter(http myhttp.Doer, host string, scheme string) *Writer {
	return &Writer{http: http, host: host, scheme: scheme}
}

func (w *Writer) StoreLesson(lesson models.Lesson) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	_, err = w.performRequestTo("POST", StoreLesson, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (w *Writer) DeleteLesson(lesson models.Lesson) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	_, err = w.performRequestTo("POST", DeleteLesson, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (w *Writer) StoreExercise(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	_, err = w.performRequestTo("POST", StoreExercise, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (w *Writer) DeleteExercise(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	_, err = w.performRequestTo("POST", DeleteExercise, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (w *Writer) IncrementBadAnswers(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	_, err = w.performRequestTo("POST", IncrementBadAnswers, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (w *Writer) IncrementGoodAnswers(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	_, err = w.performRequestTo("POST", IncrementGoodAnswers, nil, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (r *Writer) performRequestTo(method, route string, params map[string]string, reqBody []byte) ([]byte, error) {
	// parse url
	u, _ := url.Parse(r.scheme + "://" + r.host + route)

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
	resp, err := r.http.Do(req)
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
