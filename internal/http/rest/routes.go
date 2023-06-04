package rest

import (
	handlers "github.com/rtrzebinski/simple-memorizer-4/internal/http/rest/handlers"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"net/http"
)

const (
	// read

	FetchAllLessons             = "/fetch-all-lessons"
	HydrateLesson               = "/hydrate-lesson"
	FetchExercisesOfLesson      = "/fetch-exercises-of-lesson"
	FetchRandomExerciseOfLesson = "/fetch-random-exercise-of-lesson"

	// write

	StoreLesson          = "/store-lesson"
	DeleteLesson         = "/delete-lesson"
	StoreExercise        = "/store-exercise"
	DeleteExercise       = "/delete-exercise"
	IncrementBadAnswers  = "/increment-bad-answers"
	IncrementGoodAnswers = "/increment-good-answers"
)

func ListenAndServe(r storage.Reader, w storage.Writer, port string) error {
	// read

	http.Handle(FetchAllLessons, handlers.NewFetchAllLessons(r))
	http.Handle(HydrateLesson, handlers.NewHydrateLesson(r))
	http.Handle(FetchExercisesOfLesson, handlers.NewFetchExercisesOfLesson(r))
	http.Handle(FetchRandomExerciseOfLesson, handlers.NewFetchRandomExerciseOfLesson(r))

	// write

	http.Handle(StoreLesson, handlers.NewStoreLesson(w))
	http.Handle(DeleteLesson, handlers.NewDeleteLesson(w))
	http.Handle(StoreExercise, handlers.NewStoreExercise(w))
	http.Handle(DeleteExercise, handlers.NewDeleteExercise(w))
	http.Handle(IncrementBadAnswers, handlers.NewIncrementBadAnswers(w))
	http.Handle(IncrementGoodAnswers, handlers.NewIncrementGoodAnswers(w))

	if err := http.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}
