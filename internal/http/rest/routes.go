package rest

import (
	routes "github.com/rtrzebinski/simple-memorizer-4/internal/http/rest/routes"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"net/http"
)

const (
	// read

	FetchAllLessons             = "/fetch-all-lessons"
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

	http.Handle(FetchAllLessons, routes.NewFetchAllLessons(r))
	http.Handle(FetchExercisesOfLesson, routes.NewFetchExercisesOfLesson(r))
	http.Handle(FetchRandomExerciseOfLesson, routes.NewFetchRandomExerciseOfLesson(r))

	// write

	http.Handle(StoreLesson, routes.NewStoreLesson(w))
	http.Handle(DeleteLesson, routes.NewDeleteLesson(w))
	http.Handle(StoreExercise, routes.NewStoreExercise(w))
	http.Handle(DeleteExercise, routes.NewDeleteExercise(w))
	http.Handle(IncrementBadAnswers, routes.NewIncrementBadAnswers(w))
	http.Handle(IncrementGoodAnswers, routes.NewIncrementGoodAnswers(w))

	if err := http.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}
