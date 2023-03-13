package api

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-go/internal/models"
	"net/http"
)

type Client struct {
	host   string
	scheme string
}

func (c *Client) SetHost(host string) {
	c.host = host
}

func (c *Client) SetScheme(scheme string) {
	c.scheme = scheme
}

func (c *Client) GetRandomExercise() models.Exercise {
	resp, err := http.Get(c.scheme + "://" + c.host + RandomExercise)
	if err != nil {
		panic(err)
	}

	var exercise models.Exercise
	if err := json.NewDecoder(resp.Body).Decode(&exercise); err != nil {
		panic(err)
	}

	return exercise
}
