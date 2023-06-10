package models

type (
	Lesson struct {
		Id            int
		Name          string
		Description   string
		ExerciseCount int
	}

	Lessons []Lesson

	Exercise struct {
		Id          int
		Lesson      *Lesson
		Question    string
		Answer      string
		BadAnswers  int
		GoodAnswers int
	}

	Exercises []Exercise
)
