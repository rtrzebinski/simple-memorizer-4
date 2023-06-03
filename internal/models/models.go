package models

type Lesson struct {
	Id            int
	Name          string
	ExerciseCount int
}

type Lessons []Lesson

type Exercise struct {
	Id          int
	Lesson      *Lesson
	Question    string
	Answer      string
	BadAnswers  int
	GoodAnswers int
}

type Exercises []Exercise
