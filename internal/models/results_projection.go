package models

import "time"

type ResultsProjection struct {
	BadAnswers       int
	BadAnswersToday  int
	LatestBadAnswer  time.Time
	GoodAnswers      int
	GoodAnswersToday int
	LatestGoodAnswer time.Time
}

func (p ResultsProjection) GoodAnswersPercent() int {
	total := p.GoodAnswers + p.BadAnswers

	if total == 0 {
		return 0
	}

	return 100 * p.GoodAnswers / total
}
