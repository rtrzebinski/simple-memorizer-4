package validators

import "errors"

var (
	ErrLessonIdRequired         = errors.New("lesson.id is required")
	ErrLessonNameRequired       = errors.New("lesson.name is required")
	ErrLessonNameUnique         = errors.New("lesson.name must be unique")
	ErrExerciseIdRequired       = errors.New("exercise.id is required")
	ErrExerciseQuestionRequired = errors.New("exercise.question is required")
	ErrExerciseQuestionUnique   = errors.New("exercise.question must be unique")
	ErrExerciseAnswerRequired   = errors.New("exercise.answer is required")
)
