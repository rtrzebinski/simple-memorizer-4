package models

type Exercise struct {
	Id          int
	Lesson      *Lesson
	Question    string
	Answer      string
	BadAnswers  int
	GoodAnswers int
}

func (e Exercise) GoodAnswersPercent() int {
	total := e.GoodAnswers + e.BadAnswers

	if total == 0 {
		return 0
	}

	return 100 * e.GoodAnswers / total
}
