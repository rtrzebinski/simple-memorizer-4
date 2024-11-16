package result

import (
	"time"

	"github.com/rtrzebinski/simple-memorizer-4/internal/worker"
)

type ProjectionBuilder struct{}

func NewProjectionBuilder() *ProjectionBuilder {
	return &ProjectionBuilder{}
}

func (ProjectionBuilder) Projection(results []worker.Result) worker.ResultsProjection {
	projection := worker.ResultsProjection{}

	for _, a := range results {
		switch a.Type {
		case worker.Good:
			projection.GoodAnswers++
			if isToday(a.CreatedAt) {
				projection.GoodAnswersToday++
				projection.LatestGoodAnswerWasToday = true
			}
			if a.CreatedAt.After(projection.LatestGoodAnswer.Time) {
				projection.LatestGoodAnswer.Time = a.CreatedAt
			}
		case worker.Bad:
			projection.BadAnswers++
			if isToday(a.CreatedAt) {
				projection.BadAnswersToday++
				projection.LatestBadAnswerWasToday = true
			}
			if a.CreatedAt.After(projection.LatestBadAnswer.Time) {
				projection.LatestBadAnswer.Time = a.CreatedAt
			}
		}
	}

	return projection
}

func isToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}
