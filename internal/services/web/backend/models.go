package backend

import (
	"time"

	"github.com/guregu/null/v5"
)

const (
	Good ResultType = "good"
	Bad  ResultType = "bad"
)

type ResultType string

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

type Exercises []Exercise

type Lesson struct {
	Id            int
	Name          string
	Description   string
	ExerciseCount int
}

type Lessons []Lesson

type Result struct {
	Id        int
	Exercise  *Exercise
	Type      ResultType
	CreatedAt time.Time
}

type RegisterRequest struct {
	Name     string
	Email    string
	Password string
}

type RegisterResponse struct {
	AccessToken string
}

type SignInRequest struct {
	Email    string
	Password string
}

type SignInResponse struct {
	AccessToken string
}
