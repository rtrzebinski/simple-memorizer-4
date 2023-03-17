package client

import (
	"bytes"
	"encoding/json"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/server"
	"github.com/rtrzebinski/simple-memorizer-4/internal/server/routes"
	"log"
	"net/http"
	"net/url"
)

type Api struct {
	host   string
	scheme string
}

func (c *Api) Configure(url *url.URL) {
	c.host = url.Host
	c.scheme = url.Scheme
}

func (c *Api) FetchRandomExercise() models.Exercise {
	resp, err := http.Get(c.scheme + "://" + c.host + server.FetchRandomExercise)
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

func (c *Api) IncrementBadAnswers(exerciseId int) {
	input := routes.IncrementBadAnswersReq{
		ExerciseId: exerciseId,
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	urlAddress := c.scheme + "://" + c.host + server.IncrementBadAnswers
	contentType := "application/json"
	buffer := bytes.NewBuffer(jsonData)

	resp, err := http.Post(urlAddress, contentType, buffer)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(resp.StatusCode)
	}

	log.Println(resp.StatusCode)
}

func (c *Api) IncrementGoodAnswers(exerciseId int) {
	input := routes.IncrementGoodAnswersReq{
		ExerciseId: exerciseId,
	}

	jsonData, err := json.Marshal(input)
	if err != nil {
		panic(err)
	}

	urlAddress := c.scheme + "://" + c.host + server.IncrementGoodAnswers
	contentType := "application/json"
	buffer := bytes.NewBuffer(jsonData)

	resp, err := http.Post(urlAddress, contentType, buffer)

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		panic(resp.StatusCode)
	}

	log.Println(resp.StatusCode)
}
