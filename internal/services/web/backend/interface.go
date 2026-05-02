package backend

import "context"

type Reader interface {
	FetchLessons(ctx context.Context, userID string) (Lessons, error)
	HydrateLesson(ctx context.Context, userID string, lesson *Lesson) error
	FetchExercises(ctx context.Context, userID string, lesson Lesson, oldestExerciseID int) (Exercises, error)
}

type Writer interface {
	UpsertLesson(ctx context.Context, userID string, lesson *Lesson) error
	DeleteLesson(ctx context.Context, userID string, lesson Lesson) error
	UpsertExercise(ctx context.Context, userID string, exercise *Exercise) error
	StoreExercises(ctx context.Context, userID string, exercise Exercises) error
	DeleteExercise(ctx context.Context, userID string, exercise Exercise) error
}

type Publisher interface {
	PublishGoodAnswer(ctx context.Context, userID string, exerciseID int) error
	PublishBadAnswer(ctx context.Context, userID string, exerciseID int) error
}

type AuthClient interface {
	Register(ctx context.Context, firstName, lastName, email, password string) (Tokens, error)
	SignIn(ctx context.Context, email, password string) (Tokens, error)
	Revoke(ctx context.Context, refreshToken string) error
}
