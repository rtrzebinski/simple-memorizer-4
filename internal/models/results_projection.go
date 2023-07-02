package models

import "time"

type ResultsProjection struct {
	BadAnswers               int
	BadAnswersToday          int
	LatestBadAnswer          time.Time
	LatestBadAnswerWasToday  bool
	GoodAnswers              int
	GoodAnswersToday         int
	LatestGoodAnswer         time.Time
	LatestGoodAnswerWasToday bool
}

func (p *ResultsProjection) GoodAnswersPercent() int {
	total := p.GoodAnswers + p.BadAnswers

	if total == 0 {
		return 0
	}

	return 100 * p.GoodAnswers / total
}

func (p *ResultsProjection) RegisterGoodAnswer() {
	p.GoodAnswers++
	p.GoodAnswersToday++
	p.LatestGoodAnswerWasToday = true
	p.LatestGoodAnswer = time.Now()
}

func (p *ResultsProjection) RegisterBadAnswer() {
	p.BadAnswers++
	p.BadAnswersToday++
	p.LatestBadAnswerWasToday = true
	p.LatestBadAnswer = time.Now()
}
