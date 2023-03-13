package api

import (
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-go/internal/models"
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
