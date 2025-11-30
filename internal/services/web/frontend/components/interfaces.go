package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
)

// APIClient is an interface that defines methods for interacting with the backend API.
type APIClient interface {
	FetchLessons(ctx app.Context) ([]frontend.Lesson, error)
	HydrateLesson(ctx app.Context, lesson *frontend.Lesson) error
	FetchExercises(ctx app.Context, lesson frontend.Lesson, oldestExerciseID int) ([]frontend.Exercise, error)
	UpsertLesson(ctx app.Context, lesson frontend.Lesson) error
	DeleteLesson(ctx app.Context, lesson frontend.Lesson) error
	UpsertExercise(ctx app.Context, exercise frontend.Exercise) error
	StoreExercises(ctx app.Context, exercises []frontend.Exercise) error
	DeleteExercise(ctx app.Context, exercise frontend.Exercise) error
	StoreResult(ctx app.Context, result frontend.Result) error
	AuthRegister(ctx app.Context, req frontend.RegisterRequest) error
	AuthSignIn(ctx app.Context, req frontend.SignInRequest) error
	AuthLogout(ctx app.Context) error
	UserProfile(ctx app.Context) (*frontend.UserProfile, error)
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
