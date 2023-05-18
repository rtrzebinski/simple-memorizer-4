package frontend

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type ApiClient struct {
	client HttpClient
	host   string
	scheme string
}

func NewApiClient(client HttpClient, host string, scheme string) *ApiClient {
	return &ApiClient{client: client, host: host, scheme: scheme}
}

func (c *ApiClient) DeleteExercise(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := c.scheme + "://" + c.host + backend.DeleteExercise
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (c *ApiClient) StoreExercise(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := c.scheme + "://" + c.host + backend.StoreExercise
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (c *ApiClient) FetchExercisesOfLesson(lessonId int) (models.Exercises, error) {
	var output models.Exercises

	u, _ := url.Parse(c.scheme + "://" + c.host + backend.FetchExercisesOfLesson)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(lessonId))
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

func (c *ApiClient) DeleteLesson(lesson models.Lesson) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := c.scheme + "://" + c.host + backend.DeleteLesson
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (c *ApiClient) StoreLesson(lesson models.Lesson) error {
	body, err := json.Marshal(lesson)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := c.scheme + "://" + c.host + backend.StoreLesson
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (c *ApiClient) FetchAllLessons() (models.Lessons, error) {
	var output models.Lessons

	resp, err := c.performRequestTo("GET", c.scheme+"://"+c.host+backend.FetchAllLessons, nil)
	if err != nil {
		return output, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return output, fmt.Errorf("failed to decode output: %w", err)
	}

	return output, nil
}

func (c *ApiClient) FetchNextExercise() (models.Exercise, error) {
	var output models.Exercise

	resp, err := c.performRequestTo("GET", c.scheme+"://"+c.host+backend.FetchNextExercise, nil)
	if err != nil {
		return output, fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&output); err != nil {
		return output, fmt.Errorf("failed to decode output: %w", err)
	}

	return output, nil
}

func (c *ApiClient) IncrementBadAnswers(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := c.scheme + "://" + c.host + backend.IncrementBadAnswers
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (c *ApiClient) IncrementGoodAnswers(exercise models.Exercise) error {
	body, err := json.Marshal(exercise)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	u := c.scheme + "://" + c.host + backend.IncrementGoodAnswers
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", u, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (c *ApiClient) performRequestTo(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call HTTP endpoint: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		payload, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("server returned with the status code '%d': %w", resp.StatusCode, errors.New(string(payload)))
	}

	return resp, nil
}
