package memorizer

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
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
	// update previous if provided
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

		p := points(e.ResultsProjection)

		app.Log(e.Question, p)

		// populate candidates with multiplied exercise.id depending on points number
		// an exercise with 5 points has 5 times more chances to win than an exercise with 1 point etc.
		for i := 1; i <= p; i++ {
			candidates = append(candidates, e.Id)
		}
	}

	app.Log()

	// get a random winner from candidates
	winner := candidates[s.r.Intn(len(candidates))]

	// return exercise that corresponds to the winner
	return s.exercises[winner]
}

// converts results projection to points
func points(p models.ResultsProjection) int {
	// only good answers today
	if p.LatestGoodAnswerWasToday && !p.LatestBadAnswerWasToday {
		return 1
	}

	// good and bad answers today, good answer most recent
	if p.LatestGoodAnswerWasToday && p.LatestBadAnswerWasToday && p.LatestGoodAnswer.After(p.LatestBadAnswer) {
		return 1
	}

	// bad answers today
	if p.LatestBadAnswerWasToday {
		// decrease points with increasing bad answers
		if p.BadAnswersToday == 1 {
			return 80
		}

		if p.BadAnswersToday == 2 {
			return 60
		}

		if p.BadAnswersToday == 3 {
			return 40
		}

		if p.BadAnswersToday == 4 {
			return 20
		}

		return 1
	}

	// only good answers before today
	if !p.LatestGoodAnswerWasToday && !p.LatestBadAnswerWasToday && p.BadAnswers == 0 && p.GoodAnswers > 0 {
		// decrease points with increasing good answers
		if p.GoodAnswers == 1 {
			return 80
		}

		if p.GoodAnswers == 2 {
			return 60
		}

		if p.GoodAnswers == 3 {
			return 40
		}

		if p.GoodAnswers == 4 {
			return 20
		}

		return 1
	}

	// other cases

	percent := p.GoodAnswersPercent()

	if percent == 100 {
		return 1
	}

	return 100 - percent
}
