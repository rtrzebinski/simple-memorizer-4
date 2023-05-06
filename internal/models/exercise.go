package models

type Exercise struct {
	Id          int
	Question    string
	Answer      string
	BadAnswers  int
	GoodAnswers int
}

type Exercises []Exercise

type Lesson struct {
	Id   int
	Name string
}

type Lessons []Lesson
