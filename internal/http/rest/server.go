package rest

import (
	routes "github.com/rtrzebinski/simple-memorizer-4/internal/http/rest/routes"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"net/http"
)

const (
	DeleteExercise              = "/delete-exercise"
	StoreExercise               = "/store-exercise"
	FetchExercisesOfLesson      = "/fetch-exercises-of-lesson"
	DeleteLesson                = "/delete-lesson"
	StoreLesson                 = "/store-lesson"
	FetchAllLessons             = "/fetch-all-lessons"
	FetchRandomExerciseOfLesson = "/fetch-random-exercise-of-lesson"
	IncrementBadAnswers         = "/increment-bad-answers"
	IncrementGoodAnswers        = "/increment-good-answers"
)

func ListenAndServe(r storage.Reader, w storage.Writer, port string) error {
	http.Handle(DeleteExercise, routes.NewDeleteExercise(w))
	http.Handle(StoreExercise, routes.NewStoreExercise(w))
	http.Handle(FetchExercisesOfLesson, routes.NewFetchExercisesOfLesson(r))
	http.Handle(DeleteLesson, routes.NewDeleteLesson(w))
	http.Handle(StoreLesson, routes.NewStoreLesson(w))
	http.Handle(FetchAllLessons, routes.NewFetchAllLessons(r))
	http.Handle(FetchRandomExerciseOfLesson, routes.NewFetchRandomExerciseOfLesson(r))
	http.Handle(IncrementBadAnswers, routes.NewIncrementBadAnswers(w))
	http.Handle(IncrementGoodAnswers, routes.NewIncrementGoodAnswers(w))

	if err := http.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}
