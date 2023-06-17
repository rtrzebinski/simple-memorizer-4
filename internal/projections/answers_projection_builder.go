package projections

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"time"
)

func BuildAnswersProjection(answers models.Answers) models.AnswersProjection {
	projection := models.AnswersProjection{}

	for _, a := range answers {
		switch a.Type {
		case models.Good:
			projection.GoodAnswers++
			if a.CreatedAt.After(projection.LatestGoodAnswer) {
				projection.LatestGoodAnswer = a.CreatedAt
			}
			if isToday(a.CreatedAt) {
				projection.GoodAnswersToday++
			}
		case models.Bad:
			projection.BadAnswers++
			if a.CreatedAt.After(projection.LatestBadAnswer) {
				projection.LatestBadAnswer = a.CreatedAt
			}
			if isToday(a.CreatedAt) {
				projection.BadAnswersToday++
			}
		}
	}

	return projection
}

func isToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}
