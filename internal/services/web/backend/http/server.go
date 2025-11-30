package http

import (
	"net/http"
)

func ListenAndServe(s Service, v TokenVerifier, rfr TokenRefresher, port string, secure bool) error {
	// read
	http.Handle(FetchLessons, Auth(v, rfr, secure)(NewFetchLessonsHandler(s)))
	http.Handle(HydrateLesson, Auth(v, rfr, secure)(NewHydrateLessonHandler(s)))
	http.Handle(FetchExercises, Auth(v, rfr, secure)(NewFetchExercisesOfLessonHandler(s)))
	http.Handle(UserProfile, Auth(v, rfr, secure)(NewUserProfileHandler(v)))

	// write
	http.Handle(UpsertLesson, Auth(v, rfr, secure)(NewUpsertLessonHandler(s)))
	http.Handle(DeleteLesson, Auth(v, rfr, secure)(NewDeleteLessonHandler(s)))
	http.Handle(UpsertExercise, Auth(v, rfr, secure)(NewUpsertExerciseHandler(s)))
	http.Handle(StoreExercises, Auth(v, rfr, secure)(NewStoreExercisesHandler(s)))
	http.Handle(DeleteExercise, Auth(v, rfr, secure)(NewDeleteExerciseHandler(s)))
	http.Handle(StoreResult, Auth(v, rfr, secure)(NewStoreResultHandler(s)))

	// download
	http.Handle(ExportLessonCsv, Auth(v, rfr, secure)(NewExportLessonCsvHandler(s)))

	// auth
	http.Handle(AuthRegister, NewAuthRegisterHandler(s, secure))
	http.Handle(AuthSignIn, NewAuthSignInHandler(s, secure))
	http.Handle(AuthLogout, NewAuthLogoutHandler(s))

	handler := CSRFDynamicHost()(http.DefaultServeMux)

	return http.ListenAndServe(port, handler)
}
