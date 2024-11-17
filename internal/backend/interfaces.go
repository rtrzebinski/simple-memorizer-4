package backend

import "context"

type Reader interface {
	FetchLessons() (Lessons, error)
	HydrateLesson(*Lesson) error
	FetchExercises(Lesson) (Exercises, error)
}

type Writer interface {
	UpsertLesson(*Lesson) error
	DeleteLesson(Lesson) error
	UpsertExercise(*Exercise) error
	StoreExercises(Exercises) error
	DeleteExercise(Exercise) error
}

type Publisher interface {
	PublishGoodAnswer(ctx context.Context, exerciseID int) error
	PublishBadAnswer(ctx context.Context, exerciseID int) error
}
