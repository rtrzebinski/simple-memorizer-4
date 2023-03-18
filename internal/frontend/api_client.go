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

func (c *ApiClient) performRequestTo(method, url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Add("content-type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call API endpoint: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		payload, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf("API returned with status code '%d': %w", resp.StatusCode, errors.New(string(payload)))
	}

	return resp, nil
}

func (c *ApiClient) FetchRandomExercise() models.Exercise {
	resp, err := c.performRequestTo("GET", c.scheme+"://"+c.host+backend.FetchRandomExercise, nil)
	if err != nil {
		panic(fmt.Errorf("insert document HTTP error: %w", err))
	}

	defer resp.Body.Close()

	var exercise models.Exercise
	if err := json.NewDecoder(resp.Body).Decode(&exercise); err != nil {
		panic(err)
	}

	return exercise
}

func (c *ApiClient) IncrementBadAnswers(exerciseId int) {
	input := routes.IncrementBadAnswersReq{
		ExerciseId: exerciseId,
	}

	body, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	urlAddress := c.scheme + "://" + c.host + backend.IncrementBadAnswers
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", urlAddress, buffer)
	if err != nil {
		panic(fmt.Errorf("insert document HTTP error: %w", err))
	}

	defer resp.Body.Close()
}

func (c *ApiClient) IncrementGoodAnswers(exerciseId int) {
	input := routes.IncrementGoodAnswersReq{
		ExerciseId: exerciseId,
	}

	body, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	urlAddress := c.scheme + "://" + c.host + backend.IncrementGoodAnswers
	buffer := bytes.NewBuffer(body)

	resp, err := c.performRequestTo("POST", urlAddress, buffer)
	if err != nil {
		panic(fmt.Errorf("insert document HTTP error: %w", err))
	}

	defer resp.Body.Close()
}
