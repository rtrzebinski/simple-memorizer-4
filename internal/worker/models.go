package worker

import "time"

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
	LatestBadAnswer          time.Time
	LatestBadAnswerWasToday  bool
	GoodAnswers              int
	GoodAnswersToday         int
	LatestGoodAnswer         time.Time
	LatestGoodAnswerWasToday bool
}
