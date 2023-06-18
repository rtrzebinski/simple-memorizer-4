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

// Init the service with exercises to memorize
func (s *Service) Init(exercises map[int]models.Exercise) {
	s.r = rand.New(rand.NewSource(time.Now().Unix()))
	s.exercises = exercises
}

// Next exercise to memorize
func (s *Service) Next(previous models.Exercise) models.Exercise {
	// store previous back if provided
	if previous.Id > 0 {
		s.exercises[previous.Id] = previous
	}

	// holds ids of exercises
	var candidates []int

	for _, e := range s.exercises {
		// filter out previous exercise, it should never be next
		if e.Id == previous.Id {
			continue
		}
		// populate candidates with multiplied exercise.id depending on points number
		// an exercise with 5 points has 5 times more chances to win than an exercise with 1 point etc.
		for i := 1; i <= points(e.ResultsProjection.GoodAnswersPercent()); i++ {
			candidates = append(candidates, e.Id)
		}
	}

	// get a random winner from candidates
	winner := candidates[s.r.Intn(len(candidates))]

	// return exercise that corresponds to the winner
	return s.exercises[winner]
}

// converts good answers percent to points
func points(percent int) int {
	if percent <= 100 && percent > 90 {
		return 1
	} else if percent <= 90 && percent > 80 {
		return 2
	} else if percent <= 80 && percent > 70 {
		return 3
	} else if percent <= 70 && percent > 60 {
		return 4
	} else if percent <= 60 && percent > 50 {
		return 5
	} else if percent <= 50 && percent > 40 {
		return 6
	} else if percent <= 40 && percent > 30 {
		return 7
	} else if percent <= 30 && percent > 20 {
		return 8
	} else if percent <= 20 && percent > 10 {
		return 9
	} else if percent <= 10 {
		return 10
	}

	panic("Percent of good answers must be a value between 0 and 100")
}
