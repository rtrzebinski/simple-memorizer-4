package models

import (
	"github.com/guregu/null/v5"
	"time"
)

type Exercise struct {
	Id                       int
	Lesson                   *Lesson
	Question                 string
	Answer                   string
	BadAnswers               int
	BadAnswersToday          int
	LatestBadAnswer          null.Time
	LatestBadAnswerWasToday  bool
	GoodAnswers              int
	GoodAnswersToday         int
	LatestGoodAnswer         null.Time
	LatestGoodAnswerWasToday bool
}

func (e *Exercise) GoodAnswersPercent() int {
	total := e.GoodAnswers + e.BadAnswers

	if total == 0 {
		return 0
	}

	return 100 * e.GoodAnswers / total
}

func (e *Exercise) RegisterGoodAnswer() {
	e.GoodAnswers++
	e.GoodAnswersToday++
	e.LatestGoodAnswerWasToday = true
	e.LatestGoodAnswer = null.TimeFrom(time.Now())
}

func (e *Exercise) RegisterBadAnswer() {
	e.BadAnswers++
	e.BadAnswersToday++
	e.LatestBadAnswerWasToday = true
	e.LatestBadAnswer = null.TimeFrom(time.Now())
}
