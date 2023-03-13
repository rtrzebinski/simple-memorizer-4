package api

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-go/internal/models"
	"log"
	"net/http"
)

type Client struct {
	host string
}

func (c *Client) SetHost(host string) {
	c.host = host
}

func (c *Client) FetchExercise() models.Exercise {
	log.Println("Fetching exercise from the API..")

	resp, err := http.Get("http://" + c.host + Exercises)
	if err != nil {
		panic(err)
	}

	var exercise models.Exercise
	if err := json.NewDecoder(resp.Body).Decode(&exercise); err != nil {
		panic(err)
	}

	return exercise
}
