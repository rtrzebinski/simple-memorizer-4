package memorizer

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"math/rand"
	"time"
)

type Service struct {
	r         *rand.Rand
	exercises map[int]models.Exercise
}

func (s *Service) Init(exercises map[int]models.Exercise) {
	s.r = rand.New(rand.NewSource(time.Now().Unix()))
	s.exercises = exercises
}

func (s *Service) Next(previous models.Exercise) models.Exercise {
	// will update number of answers of previous if provided
	if previous.Id > 0 {
		s.exercises[previous.Id] = previous
	}

	var exercises []models.Exercise

	// convert into slice and filter out previous exercise
	for _, e := range s.exercises {
		if e.Id != previous.Id {
			exercises = append(exercises, e)
		}
	}

	// return random
	return exercises[s.r.Intn(len(exercises))]
}
