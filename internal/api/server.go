package api

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/api/methods"
	"github.com/rtrzebinski/simple-memorizer-4/internal/storage"
	"net/http"
)

const (
	FetchRandomExercise  = "/fetch-random-exercise"
	IncrementBadAnswers  = "/increment-bad-answers"
	IncrementGoodAnswers = "/increment-good-answers"
)

func ListenAndServe(r storage.Reader, w storage.Writer, port string) {
	http.Handle(FetchRandomExercise, methods.NewFetchRandomExercise(r))
	http.Handle(IncrementBadAnswers, methods.NewIncrementBadAnswers(w))
	http.Handle(IncrementGoodAnswers, methods.NewIncrementGoodAnswers(w))

	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
