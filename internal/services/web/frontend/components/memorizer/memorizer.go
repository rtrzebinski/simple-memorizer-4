package memorizer

import (
	"math/rand"
	"time"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
)

// Memorizer provides next exercise to memorize
// points are calculated based on exercise results
// the more points an exercise has the more chances it has to be picked
type Memorizer struct {
	r         *rand.Rand
	exercises map[int]frontend.Exercise
}

// Init the service with exercises to memorize
func (s *Memorizer) Init(exercises map[int]frontend.Exercise) {
	s.r = rand.New(rand.NewSource(time.Now().Unix()))
	s.exercises = exercises
}

// Next exercise to memorize
func (s *Memorizer) Next(previous frontend.Exercise) frontend.Exercise {
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

		p := points(e)

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
func points(e frontend.Exercise) int {
	// only good answers today
	if e.LatestGoodAnswerWasToday && !e.LatestBadAnswerWasToday {
		return 1
	}

	// good and bad answers today, good answer most recent
	if e.LatestGoodAnswerWasToday && e.LatestBadAnswerWasToday && e.LatestGoodAnswer.Time.After(e.LatestBadAnswer.Time) {
		return 1
	}

	// bad answers today
	if e.LatestBadAnswerWasToday {
		// decrease points with increasing bad answers
		if e.BadAnswersToday == 1 {
			return 80
		}

		if e.BadAnswersToday == 2 {
			return 60
		}

		if e.BadAnswersToday == 3 {
			return 40
		}

		if e.BadAnswersToday == 4 {
			return 20
		}

		return 1
	}

	// only good answers before today
	if !e.LatestGoodAnswerWasToday && !e.LatestBadAnswerWasToday && e.BadAnswers == 0 && e.GoodAnswers > 0 {
		// decrease points with increasing good answers
		if e.GoodAnswers == 1 {
			return 80
		}

		if e.GoodAnswers == 2 {
			return 60
		}

		if e.GoodAnswers == 3 {
			return 40
		}

		if e.GoodAnswers == 4 {
			return 20
		}

		return 1
	}

	// other cases

	percent := e.GoodAnswersPercent()

	if percent == 100 {
		return 1
	}

	return 100 - percent
}
