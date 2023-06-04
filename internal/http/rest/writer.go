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

	u := w.scheme + "://" + w.host + StoreLesson

	_, err = w.performRequestTo("POST", u, body)
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

	u := w.scheme + "://" + w.host + DeleteLesson

	_, err = w.performRequestTo("POST", u, body)
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

	u := w.scheme + "://" + w.host + StoreExercise

	_, err = w.performRequestTo("POST", u, body)
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

	u := w.scheme + "://" + w.host + DeleteExercise

	_, err = w.performRequestTo("POST", u, body)
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

	u := w.scheme + "://" + w.host + IncrementBadAnswers

	_, err = w.performRequestTo("POST", u, body)
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

	u := w.scheme + "://" + w.host + IncrementGoodAnswers

	_, err = w.performRequestTo("POST", u, body)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	return nil
}

func (r *Writer) performRequestTo(method, url string, reqBody []byte) ([]byte, error) {
	buffer := bytes.NewBuffer(reqBody)

	req, err := http.NewRequest(method, url, buffer)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	resp, err := r.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call HTTP endpoint: %w", err)
	}

	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read HTTP response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("server returned with the status code '%d': %w",
			resp.StatusCode, errors.New(string(respBody)))
	}

	return respBody, nil
}
