package backend

import "context"

type Reader interface {
	FetchLessons(ctx context.Context) (Lessons, error)
	HydrateLesson(ctx context.Context, lesson *Lesson) error
	FetchExercises(ctx context.Context, lesson Lesson) (Exercises, error)
}

type Writer interface {
	UpsertLesson(ctx context.Context, lesson *Lesson) error
	DeleteLesson(ctx context.Context, lesson Lesson) error
	UpsertExercise(ctx context.Context, exercise *Exercise) error
	StoreExercises(ctx context.Context, exercise Exercises) error
	DeleteExercise(ctx context.Context, exercise Exercise) error
}

type Publisher interface {
	PublishGoodAnswer(ctx context.Context, exerciseID int) error
	PublishBadAnswer(ctx context.Context, exerciseID int) error
}
