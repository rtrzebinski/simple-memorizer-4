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

	respBody, err := r.performRequestTo("GET", FetchAllLessons, nil, nil)
	if err != nil {
		return lessons, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &lessons); err != nil {
		return lessons, fmt.Errorf("failed to decode lessons: %w", err)
	}

	return lessons, nil
}

func (r *Reader) HydrateLesson(lesson *models.Lesson) error {
	var params = map[string]string{
		"lesson_id": strconv.Itoa(lesson.Id),
	}

	respBody, err := r.performRequestTo("GET", HydrateLesson, params, nil)
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

	var params = map[string]string{
		"lesson_id": strconv.Itoa(lesson.Id),
	}

	respBody, err := r.performRequestTo("GET", FetchExercisesOfLesson, params, nil)
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

	var params = map[string]string{
		"lesson_id": strconv.Itoa(lesson.Id),
	}

	respBody, err := r.performRequestTo("GET", FetchRandomExerciseOfLesson, params, nil)
	if err != nil {
		return exercise, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	if err := json.Unmarshal(respBody, &exercise); err != nil {
		return exercise, fmt.Errorf("failed to decode exercise: %w", err)
	}

	return exercise, nil
}

func (r *Reader) performRequestTo(method, route string, params map[string]string, reqBody []byte) ([]byte, error) {
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
