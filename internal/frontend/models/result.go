package models

type ResultType string

const (
	Good ResultType = "good"
	Bad  ResultType = "bad"
)

type Result struct {
	Id       int
	Exercise *Exercise
	Type     ResultType
}
