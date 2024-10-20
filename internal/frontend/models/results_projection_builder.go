package models

import (
	"time"
)

func BuildResultsProjection(results Results) ResultsProjection {
	projection := ResultsProjection{}

	for _, a := range results {
		switch a.Type {
		case Good:
			projection.GoodAnswers++
			if isToday(a.CreatedAt) {
				projection.GoodAnswersToday++
				projection.LatestGoodAnswerWasToday = true
			}
			if a.CreatedAt.After(projection.LatestGoodAnswer) {
				projection.LatestGoodAnswer = a.CreatedAt
			}
		case Bad:
			projection.BadAnswers++
			if isToday(a.CreatedAt) {
				projection.BadAnswersToday++
				projection.LatestBadAnswerWasToday = true
			}
			if a.CreatedAt.After(projection.LatestBadAnswer) {
				projection.LatestBadAnswer = a.CreatedAt
			}
		}
	}

	return projection
}

func isToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}
