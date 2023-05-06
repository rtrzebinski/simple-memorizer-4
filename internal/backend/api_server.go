package backend

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/routes"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/storage"
	"net/http"
)

const (
	DeleteExercise       = "/delete-exercise"
	StoreExercise        = "/store-exercise"
	FetchAllExercises    = "/fetch-all-exercises"
	DeleteLesson         = "/delete-lesson"
	StoreLesson          = "/store-lesson"
	FetchAllLessons      = "/fetch-all-lessons"
	FetchNextExercise    = "/fetch-next-exercise"
	IncrementBadAnswers  = "/increment-bad-answers"
	IncrementGoodAnswers = "/increment-good-answers"
)

func ListenAndServe(r storage.Reader, w storage.Writer, port string) error {
	http.Handle(DeleteExercise, routes.NewDeleteExercise(w))
	http.Handle(StoreExercise, routes.NewStoreExercise(w))
	http.Handle(FetchAllExercises, routes.NewFetchAllExercises(r))
	http.Handle(FetchNextExercise, routes.NewFetchNextExercise(r))
	http.Handle(IncrementBadAnswers, routes.NewIncrementBadAnswers(w))
	http.Handle(IncrementGoodAnswers, routes.NewIncrementGoodAnswers(w))

	if err := http.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}
