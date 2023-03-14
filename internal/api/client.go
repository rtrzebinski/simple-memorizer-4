package api

import (
	"bytes"
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-go/internal/models"
	"log"
	"net/http"
	"net/url"
)

type Client struct {
	host   string
	scheme string
}

func (c *Client) Configure(url *url.URL) {
	c.host = url.Host
	c.scheme = url.Scheme
}

func (c *Client) FetchRandomExercise() models.Exercise {
	resp, err := http.Get(c.scheme + "://" + c.host + FetchRandomExercise)
	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	var exercise models.Exercise
	if err := json.NewDecoder(resp.Body).Decode(&exercise); err != nil {
		panic(err)
	}

	return exercise
}

func (c *Client) IncrementBadAnswers(exerciseId int) {
	type Input struct {
		ExerciseId int `json:"exercise_id"`
	}

	input := Input{
		ExerciseId: exerciseId,
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	urlAddress := c.scheme + "://" + c.host + IncrementBadAnswers
	contentType := "application/json"
	buffer := bytes.NewBuffer(jsonData)

	resp, err := http.Post(urlAddress, contentType, buffer)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(resp.StatusCode)
	}

	log.Println(resp.StatusCode)
}

func (c *Client) IncrementGoodAnswers(exerciseId int) {
	type Input struct {
		ExerciseId int `json:"exercise_id"`
	}

	input := Input{
		ExerciseId: exerciseId,
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	urlAddress := c.scheme + "://" + c.host + IncrementGoodAnswers
	contentType := "application/json"
	buffer := bytes.NewBuffer(jsonData)

	resp, err := http.Post(urlAddress, contentType, buffer)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(resp.StatusCode)
	}

	log.Println(resp.StatusCode)
}
