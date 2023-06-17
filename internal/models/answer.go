package models

import "time"

type AnswerType string

const (
	Good AnswerType = "good"
	Bad  AnswerType = "bad"
)

type Answer struct {
	Id        int
	Exercise  *Exercise
	Type      AnswerType
	CreatedAt time.Time
}
