package backend

import "time"

const (
	Good ResultType = "good"
	Bad  ResultType = "bad"
)

type Exercise struct {
	Id       int
	Lesson   *Lesson
	Question string
	Answer   string
	Results  Results
}

type Exercises []Exercise

type Lesson struct {
	Id            int
	Name          string
	Description   string
	ExerciseCount int
}

type Lessons []Lesson

type ResultType string

type Result struct {
	Id        int
	Exercise  *Exercise
	Type      ResultType
	CreatedAt time.Time
}

type Results []Result
