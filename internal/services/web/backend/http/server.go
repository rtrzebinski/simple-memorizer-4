package http

import (
	"net/http"
)

func ListenAndServe(s Service, port string) error {
	// read

	http.Handle(FetchLessons, NewFetchLessonsHandler(s))
	http.Handle(HydrateLesson, NewHydrateLessonHandler(s))
	http.Handle(FetchExercises, NewFetchExercisesOfLessonHandler(s))

	// write

	http.Handle(UpsertLesson, NewUpsertLessonHandler(s))
	http.Handle(DeleteLesson, NewDeleteLessonHandler(s))
	http.Handle(UpsertExercise, NewUpsertExerciseHandler(s))
	http.Handle(StoreExercises, NewStoreExercisesHandler(s))
	http.Handle(DeleteExercise, NewDeleteExerciseHandler(s))
	http.Handle(StoreResult, NewStoreResultHandler(s))

	// download

	http.Handle(ExportLessonCsv, NewExportLessonCsvHandler(s))

	// auth

	http.Handle(AuthRegister, NewAuthRegisterHandler(s))
	http.Handle(AuthSignIn, NewAuthSignInHandler(s))

	if err := http.ListenAndServe(port, nil); err != nil {
		return err
	}

	return nil
}
