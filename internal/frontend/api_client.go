package frontend

import (
	"bytes"
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/routes"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"log"
	"net/http"
	"net/url"
)

type ApiClient struct {
	host   string
	scheme string
}

func (c *ApiClient) Configure(url *url.URL) {
	c.host = url.Host
	c.scheme = url.Scheme
}

func (c *ApiClient) FetchRandomExercise() models.Exercise {
	resp, err := http.Get(c.scheme + "://" + c.host + backend.FetchRandomExercise)
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

func (c *ApiClient) IncrementBadAnswers(exerciseId int) {
	input := routes.IncrementBadAnswersReq{
		ExerciseId: exerciseId,
	}

	body, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	urlAddress := c.scheme + "://" + c.host + backend.IncrementBadAnswers
	contentType := "application/json"
	buffer := bytes.NewBuffer(body)

	resp, err := http.Post(urlAddress, contentType, buffer)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(resp.StatusCode)
	}

	log.Println(resp.StatusCode)
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
	contentType := "application/json"
	buffer := bytes.NewBuffer(body)

	resp, err := http.Post(urlAddress, contentType, buffer)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(resp.StatusCode)
	}

	log.Println(resp.StatusCode)
}
