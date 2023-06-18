package models

type Exercise struct {
	Id                int
	Lesson            *Lesson
	Question          string
	Answer            string
	Results           Results
	ResultsProjection ResultsProjection
}
