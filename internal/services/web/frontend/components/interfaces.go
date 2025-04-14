package components

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
)

type APIClient interface {
	FetchLessons(ctx context.Context, authToken string) ([]frontend.Lesson, error)
	HydrateLesson(ctx context.Context, lesson *frontend.Lesson, authToken string) error
	FetchExercises(ctx context.Context, lesson frontend.Lesson, authToken string) ([]frontend.Exercise, error)
	UpsertLesson(ctx context.Context, lesson frontend.Lesson, authToken string) error
	DeleteLesson(ctx context.Context, lesson frontend.Lesson, authToken string) error
	UpsertExercise(ctx context.Context, exercise frontend.Exercise, authToken string) error
	StoreExercises(ctx context.Context, exercises []frontend.Exercise, authToken string) error
	DeleteExercise(ctx context.Context, exercise frontend.Exercise, authToken string) error
	StoreResult(ctx context.Context, result frontend.Result, authToken string) error
	AuthRegister(ctx context.Context, req frontend.RegisterRequest) (frontend.RegisterResponse, error)
	AuthSignIn(ctx context.Context, req frontend.SignInRequest) (frontend.SignInResponse, error)
}
