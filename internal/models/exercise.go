package models

type Exercise struct {
	Id          int
	Question    string
	Answer      string
	BadAnswers  int
	GoodAnswers int
}

type Exercises []Exercise
