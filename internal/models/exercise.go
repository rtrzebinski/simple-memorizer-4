package models

type Exercise struct {
	Id                int
	Lesson            *Lesson
	Question          string
	Answer            string
	Answers           Answers
	AnswersProjection AnswersProjection
}
