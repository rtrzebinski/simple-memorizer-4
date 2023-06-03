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

type Client struct {
	http   myhttp.Doer
	host   string
	scheme string
}

func NewClient(http myhttp.Doer, host string, scheme string) *Client {
	return &Client{http: http, host: host, scheme: scheme}
}

func (c *Client) DeleteExercise(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := c.scheme + "://" + c.host + DeleteExercise
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (c *Client) StoreExercise(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := c.scheme + "://" + c.host + StoreExercise
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (c *Client) FetchExercisesOfLesson(lesson models.Lesson) (models.Exercises, error) {
	var output models.Exercises

	u, _ := url.Parse(c.scheme + "://" + c.host + FetchExercisesOfLesson)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lesson.Id))
	u.RawQuery = params.Encode()

	resp, err := c.performRequestTo("GET", u.String(), nil)
	if err != nil {
		return output, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return output, fmt.Errorf("failed to decode output: %w", err)
	}

	return output, nil
}

func (c *Client) DeleteLesson(lesson models.Lesson) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := c.scheme + "://" + c.host + DeleteLesson
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (c *Client) StoreLesson(lesson models.Lesson) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := c.scheme + "://" + c.host + StoreLesson
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (c *Client) FetchAllLessons() (models.Lessons, error) {
	var output models.Lessons

	resp, err := c.performRequestTo("GET", c.scheme+"://"+c.host+FetchAllLessons, nil)
	if err != nil {
		return output, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return output, fmt.Errorf("failed to decode output: %w", err)
	}

	return output, nil
}

func (c *Client) FetchRandomExerciseOfLesson(lesson models.Lesson) (models.Exercise, error) {
	var output models.Exercise

	u, _ := url.Parse(c.scheme + "://" + c.host + FetchRandomExerciseOfLesson)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lesson.Id))
	u.RawQuery = params.Encode()

	resp, err := c.performRequestTo("GET", u.String(), nil)
	if err != nil {
		return output, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return output, fmt.Errorf("failed to decode output: %w", err)
	}

	return output, nil
}

func (c *Client) IncrementBadAnswers(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := c.scheme + "://" + c.host + IncrementBadAnswers
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (c *Client) IncrementGoodAnswers(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := c.scheme + "://" + c.host + IncrementGoodAnswers
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (c *Client) performRequestTo(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call HTTP endpoint: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		payload, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("server returned with the status code '%d': %w", resp.StatusCode, errors.New(string(payload)))
	}

	return resp, nil
}
