package frontend

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/routes"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"io"
	"net/http"
)

type ApiClient struct {
	client HttpClient
	host   string
	scheme string
}

func NewApiClient(client HttpClient, host string, scheme string) *ApiClient {
	return &ApiClient{client: client, host: host, scheme: scheme}
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

func (c *ApiClient) IncrementBadAnswers(exerciseId int) error {
	input := routes.IncrementBadAnswersReq{
		ExerciseId: exerciseId,
	}

	body, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	urlAddress := c.scheme + "://" + c.host + backend.IncrementBadAnswers
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", urlAddress, buffer)
	if err != nil {
		return fmt.Errorf("failed to perform HTTP request: %w", err)
	}

	defer resp.Body.Close()

	return nil
}

func (c *ApiClient) IncrementGoodAnswers(exerciseId int) error {
	input := routes.IncrementGoodAnswersReq{
		ExerciseId: exerciseId,
	}

	body, err := json.Marshal(input)
	if err != nil {
		return fmt.Errorf("failed to encode input: %w", err)
	}

	urlAddress := c.scheme + "://" + c.host + backend.IncrementGoodAnswers
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", urlAddress, buffer)
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
