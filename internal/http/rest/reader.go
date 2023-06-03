package rest

import (
	"encoding/json"
	"errors"
	"fmt"
	myhttp "github.com/rtrzebinski/simple-memorizer-4/internal/http"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type Reader struct {
	http   myhttp.Doer
	host   string
	scheme string
}

func NewReader(http myhttp.Doer, host string, scheme string) *Reader {
	return &Reader{http: http, host: host, scheme: scheme}
}

func (r *Reader) FetchExercisesOfLesson(lesson models.Lesson) (models.Exercises, error) {
	var output models.Exercises

	u, _ := url.Parse(r.scheme + "://" + r.host + FetchExercisesOfLesson)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lesson.Id))
	u.RawQuery = params.Encode()

	resp, err := r.performRequestTo("GET", u.String(), nil)
	if err != nil {
		return output, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return output, fmt.Errorf("failed to decode output: %w", err)
	}

	return output, nil
}

func (r *Reader) FetchAllLessons() (models.Lessons, error) {
	var output models.Lessons

	resp, err := r.performRequestTo("GET", r.scheme+"://"+r.host+FetchAllLessons, nil)
	if err != nil {
		return output, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return output, fmt.Errorf("failed to decode output: %w", err)
	}

	return output, nil
}

func (r *Reader) FetchRandomExerciseOfLesson(lesson models.Lesson) (models.Exercise, error) {
	var output models.Exercise

	u, _ := url.Parse(r.scheme + "://" + r.host + FetchRandomExerciseOfLesson)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lesson.Id))
	u.RawQuery = params.Encode()

	resp, err := r.performRequestTo("GET", u.String(), nil)
	if err != nil {
		return output, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return output, fmt.Errorf("failed to decode output: %w", err)
	}

	return output, nil
}

func (r *Reader) performRequestTo(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	resp, err := r.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call HTTP endpoint: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		payload, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("server returned with the status code '%d': %w", resp.StatusCode, errors.New(string(payload)))
	}

	return resp, nil
}
