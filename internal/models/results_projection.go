package models

import "time"

type ResultsProjection struct {
	LatestBadAnswer  time.Time
	BadAnswers       int
	BadAnswersToday  int
	LatestGoodAnswer time.Time
	GoodAnswers      int
	GoodAnswersToday int
}

func (p ResultsProjection) GoodAnswersPercent() int {
	total := p.GoodAnswers + p.BadAnswers

	if total == 0 {
		return 0
	}

	return 100 * p.GoodAnswers / total
}
