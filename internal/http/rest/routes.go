package rest

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal"
	handlers "github.com/rtrzebinski/simple-memorizer-4/internal/http/rest/handlers"
	"net/http"
)

const (
	// read

	FetchLessons   = "/fetch-lessons"
	HydrateLesson  = "/hydrate-lesson"
	FetchExercises = "/fetch-exercises"

	// write

	UpsertLesson   = "/upsert-lesson"
	DeleteLesson   = "/delete-lesson"
	UpsertExercise = "/upsert-exercise"
	DeleteExercise = "/delete-exercise"
	StoreResult    = "/store-result"

	// download

	ExportLessonCsv = "/export-lesson-csv"
)

func ListenAndServe(r internal.Reader, w internal.Writer, port string) error {
	// read

	http.Handle(FetchLessons, handlers.NewFetchLessons(r))
	http.Handle(HydrateLesson, handlers.NewHydrateLesson(r))
	http.Handle(FetchExercises, handlers.NewFetchExercisesOfLesson(r))

	// write

	http.Handle(UpsertLesson, handlers.NewUpsertLesson(w))
	http.Handle(DeleteLesson, handlers.NewDeleteLesson(w))
	http.Handle(UpsertExercise, handlers.NewUpsertExercise(w))
	http.Handle(DeleteExercise, handlers.NewDeleteExercise(w))
	http.Handle(StoreResult, handlers.NewStoreResult(w))

	// download

	http.Handle(ExportLessonCsv, handlers.NewExportLessonCsv(r))

	if err := http.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}
