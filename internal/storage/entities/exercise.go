package entities

type Exercise struct {
	Id       int
	Question string
	Answer   string
}

type ExerciseResult struct {
	Id          int
	ExerciseId  int
	BadAnswers  int
	GoodAnswers int
}
