package memorizer

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"math/rand"
	"time"
)

type Service struct {
	r         *rand.Rand
	exercises []models.Exercise
}

func (s *Service) Init(exercises []models.Exercise) {
	s.r = rand.New(rand.NewSource(time.Now().Unix()))
	s.exercises = exercises
}

func (s *Service) Next(previous models.Exercise) models.Exercise {
	var exercises []models.Exercise

	// filter out previous
	for _, e := range s.exercises {
		if e.Id != previous.Id {
			exercises = append(exercises, e)
		}
	}

	// return random
	return exercises[s.r.Intn(len(exercises))]
}
