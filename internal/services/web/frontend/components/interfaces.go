package components

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
)

type APIClient interface {
	FetchLessons(ctx context.Context, accessToken string) ([]frontend.Lesson, error)
	HydrateLesson(ctx context.Context, lesson *frontend.Lesson, accessToken string) error
	FetchExercises(ctx context.Context, lesson frontend.Lesson, oldestExerciseID int, accessToken string) ([]frontend.Exercise, error)
	UpsertLesson(ctx context.Context, lesson frontend.Lesson, accessToken string) error
	DeleteLesson(ctx context.Context, lesson frontend.Lesson, accessToken string) error
	UpsertExercise(ctx context.Context, exercise frontend.Exercise, accessToken string) error
	StoreExercises(ctx context.Context, exercises []frontend.Exercise, accessToken string) error
	DeleteExercise(ctx context.Context, exercise frontend.Exercise, accessToken string) error
	StoreResult(ctx context.Context, result frontend.Result, accessToken string) error
	AuthRegister(ctx context.Context, req frontend.RegisterRequest) (frontend.RegisterResponse, error)
	AuthSignIn(ctx context.Context, req frontend.SignInRequest) (frontend.SignInResponse, error)
}
