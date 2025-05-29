package backend

import "context"

type Reader interface {
	FetchLessons(ctx context.Context, userID string) (Lessons, error)
	HydrateLesson(ctx context.Context, lesson *Lesson, userID string) error
	FetchExercises(ctx context.Context, lesson Lesson, oldestExerciseID int, userID string) (Exercises, error)
}

type Writer interface {
	UpsertLesson(ctx context.Context, lesson *Lesson, userID string) error
	DeleteLesson(ctx context.Context, lesson Lesson, userID string) error
	UpsertExercise(ctx context.Context, exercise *Exercise, userID string) error
	StoreExercises(ctx context.Context, exercise Exercises, userID string) error
	DeleteExercise(ctx context.Context, exercise Exercise, userID string) error
}

type Publisher interface {
	PublishGoodAnswer(ctx context.Context, exerciseID int) error
	PublishBadAnswer(ctx context.Context, exerciseID int) error
}

type AuthClient interface {
	Register(ctx context.Context, name, email, password string) (accessToken string, err error)
	SignIn(ctx context.Context, email, password string) (accessToken string, err error)
}
