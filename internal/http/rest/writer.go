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
	buffer := bytes.NewBuffer(body)

	resp, err := w.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (w *Writer) DeleteLesson(lesson models.Lesson) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := w.scheme + "://" + w.host + DeleteLesson
	buffer := bytes.NewBuffer(body)

	resp, err := w.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (w *Writer) StoreExercise(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := w.scheme + "://" + w.host + StoreExercise
	buffer := bytes.NewBuffer(body)

	resp, err := w.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (w *Writer) DeleteExercise(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := w.scheme + "://" + w.host + DeleteExercise
	buffer := bytes.NewBuffer(body)

	resp, err := w.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (w *Writer) IncrementBadAnswers(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := w.scheme + "://" + w.host + IncrementBadAnswers
	buffer := bytes.NewBuffer(body)

	resp, err := w.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (w *Writer) IncrementGoodAnswers(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := w.scheme + "://" + w.host + IncrementGoodAnswers
	buffer := bytes.NewBuffer(body)

	resp, err := w.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (w *Writer) performRequestTo(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	resp, err := w.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call HTTP endpoint: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		payload, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("server returned with the status code '%d': %w", resp.StatusCode, errors.New(string(payload)))
	}

	return resp, nil
}
