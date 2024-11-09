package worker

type ResultType string

const (
	Good ResultType = "good"
	Bad  ResultType = "bad"
)

type Result struct {
	Id         int
	ExerciseId int
	Type       ResultType
}
