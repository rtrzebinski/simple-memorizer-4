package server

import (
	"net/http"
)

func ListenAndServe(r Reader, w Writer, p Publisher, port, certFile, keyFile string) error {
	// read

	http.Handle(FetchLessons, NewFetchLessonsHandler(r))
	http.Handle(HydrateLesson, NewHydrateLessonHandler(r))
	http.Handle(FetchExercises, NewFetchExercisesOfLessonHandler(r))

	// write

	http.Handle(UpsertLesson, NewUpsertLessonHandler(w))
	http.Handle(DeleteLesson, NewDeleteLessonHandler(w))
	http.Handle(UpsertExercise, NewUpsertExerciseHandler(w))
	http.Handle(StoreExercises, NewStoreExercisesHandler(w))
	http.Handle(DeleteExercise, NewDeleteExerciseHandler(w))
	http.Handle(StoreResult, NewStoreResultHandler(p))

	// download

	http.Handle(ExportLessonCsv, NewExportLessonCsvHandler(r))

	if err := http.ListenAndServeTLS(port, certFile, keyFile, nil); err != nil {
		return err
	}

	return nil
}
