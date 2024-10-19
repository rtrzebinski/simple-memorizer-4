package server

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/routes"
	"net/http"
)

func ListenAndServe(r Reader, w Writer, port, certFile, keyFile string) error {
	// read

	http.Handle(routes.FetchLessons, NewFetchLessons(r))
	http.Handle(routes.HydrateLesson, NewHydrateLesson(r))
	http.Handle(routes.FetchExercises, NewFetchExercisesOfLesson(r))

	// write

	http.Handle(routes.UpsertLesson, NewUpsertLesson(w))
	http.Handle(routes.DeleteLesson, NewDeleteLesson(w))
	http.Handle(routes.UpsertExercise, NewUpsertExercise(w))
	http.Handle(routes.StoreExercises, NewStoreExercises(w))
	http.Handle(routes.DeleteExercise, NewDeleteExercise(w))
	http.Handle(routes.StoreResult, NewStoreResult(w))

	// download

	http.Handle(routes.ExportLessonCsv, NewExportLessonCsv(r))

	if err := http.ListenAndServeTLS(port, certFile, keyFile, nil); err != nil {
		return err
	}

	return nil
}
