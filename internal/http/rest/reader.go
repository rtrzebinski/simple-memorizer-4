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

func (r *Reader) FetchAllLessons() (models.Lessons, error) {
	var lessons models.Lessons

	respBody, err := r.performRequestTo("GET", r.scheme+"://"+r.host+FetchAllLessons, nil)
	if err != nil {
		return lessons, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &lessons); err != nil {
		return lessons, fmt.Errorf("failed to decode lessons: %w", err)
	}

	return lessons, nil
}

func (r *Reader) HydrateLesson(lesson *models.Lesson) error {
	u, _ := url.Parse(r.scheme + "://" + r.host + HydrateLesson)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lesson.Id))
	u.RawQuery = params.Encode()

	respBody, err := r.performRequestTo("GET", u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, lesson); err != nil {
		return fmt.Errorf("failed to decode lesson: %w", err)
	}

	return nil
}

func (r *Reader) FetchExercisesOfLesson(lesson models.Lesson) (models.Exercises, error) {
	var exercises models.Exercises

	u, _ := url.Parse(r.scheme + "://" + r.host + FetchExercisesOfLesson)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lesson.Id))
	u.RawQuery = params.Encode()

	respBody, err := r.performRequestTo("GET", u.String(), nil)
	if err != nil {
		return exercises, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &exercises); err != nil {
		return exercises, fmt.Errorf("failed to decode exercises: %w", err)
	}

	return exercises, nil
}

func (r *Reader) FetchRandomExerciseOfLesson(lesson models.Lesson) (models.Exercise, error) {
	var exercise models.Exercise

	u, _ := url.Parse(r.scheme + "://" + r.host + FetchRandomExerciseOfLesson)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lesson.Id))
	u.RawQuery = params.Encode()

	respBody, err := r.performRequestTo("GET", u.String(), nil)
	if err != nil {
		return exercise, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &exercise); err != nil {
		return exercise, fmt.Errorf("failed to decode exercise: %w", err)
	}

	return exercise, nil
}

func (r *Reader) performRequestTo(method, url string, reqBody []byte) ([]byte, error) {
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
