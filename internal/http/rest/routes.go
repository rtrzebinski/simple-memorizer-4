package rest

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal"
	handlers "github.com/rtrzebinski/simple-memorizer-4/internal/http/rest/handlers"
	"net/http"
)

const (
	// read

	FetchAllLessons        = "/fetch-all-lessons"
	HydrateLesson          = "/hydrate-lesson"
	FetchExercisesOfLesson = "/fetch-exercises-of-lesson"

	// write

	StoreLesson    = "/store-lesson"
	DeleteLesson   = "/delete-lesson"
	StoreExercise  = "/store-exercise"
	DeleteExercise = "/delete-exercise"
	StoreResult    = "/store-result"
)

func ListenAndServe(r internal.Reader, w internal.Writer, port string) error {
	// read

	http.Handle(FetchAllLessons, handlers.NewFetchAllLessons(r))
	http.Handle(HydrateLesson, handlers.NewHydrateLesson(r))
	http.Handle(FetchExercisesOfLesson, handlers.NewFetchExercisesOfLesson(r))

	// write

	http.Handle(StoreLesson, handlers.NewStoreLesson(w))
	http.Handle(DeleteLesson, handlers.NewDeleteLesson(w))
	http.Handle(StoreExercise, handlers.NewStoreExercise(w))
	http.Handle(DeleteExercise, handlers.NewDeleteExercise(w))
	http.Handle(StoreResult, handlers.NewStoreResult(w))

	if err := http.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}
