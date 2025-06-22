package components

import (
	"context"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
)

// APIClient is an interface that defines methods for interacting with the backend API.
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

// ExerciseEditor is an interface that defines methods for components that can edit exercises.
type ExerciseEditor interface {
	// getLesson returns the lesson that edited exercise belongs to, used for FK creation.
	getLesson() *frontend.Lesson
	// getExercises returns the list of exercises of the lesson, used for validation against duplicates.
	getExercises() []*frontend.Exercise
	// exerciseEditDone is called when the exercise edit is done - either submitted or cancelled.
	exerciseEditDone(ctx app.Context)
}
