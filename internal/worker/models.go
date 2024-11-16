package worker

import (
	"time"

	"github.com/guregu/null/v5"
)

type ResultType string

const (
	Good ResultType = "good"
	Bad  ResultType = "bad"
)

type Result struct {
	Id         int
	ExerciseId int
	Type       ResultType
	CreatedAt  time.Time
}

type ResultsProjection struct {
	BadAnswers               int
	BadAnswersToday          int
	LatestBadAnswer          null.Time
	LatestBadAnswerWasToday  bool
	GoodAnswers              int
	GoodAnswersToday         int
	LatestGoodAnswer         null.Time
	LatestGoodAnswerWasToday bool
}
