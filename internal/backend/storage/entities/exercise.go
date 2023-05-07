package entities

type Exercise struct {
	Id       int
	LessonId int
	Question string
	Answer   string
}

type Lesson struct {
	Id   int
	Name string
}
