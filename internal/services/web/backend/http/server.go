package http

import (
	"net/http"
)

const (
	FetchLessons    = "/fetch-lessons"
	HydrateLesson   = "/hydrate-lesson"
	FetchExercises  = "/fetch-exercises"
	UserProfile     = "/user-profile"
	UpsertLesson    = "/upsert-lesson"
	DeleteLesson    = "/delete-lesson"
	UpsertExercise  = "/upsert-exercise"
	StoreExercises  = "/store-exercises"
	DeleteExercise  = "/delete-exercise"
	StoreResult     = "/store-result"
	ExportLessonCsv = "/export-lesson-csv"
	AuthRegister    = "/auth-register"
	AuthSignIn      = "/auth-sign-in"
	AuthLogout      = "/auth-logout"
)

func ListenAndServe(s Service, v TokenVerifier, rfr TokenRefresher, port string, secure bool) error {
	// read
	http.Handle(FetchLessons, Auth(v, rfr, secure)(NewHandlerFetchLessons(s)))
	http.Handle(HydrateLesson, Auth(v, rfr, secure)(NewHandlerHydrateLesson(s)))
	http.Handle(FetchExercises, Auth(v, rfr, secure)(NewHandlerFetchExercisesOfLesson(s)))
	http.Handle(UserProfile, Auth(v, rfr, secure)(NewUserProfileHandler(v)))

	// write
	http.Handle(UpsertLesson, Auth(v, rfr, secure)(NewHandlerUpsertLesson(s)))
	http.Handle(DeleteLesson, Auth(v, rfr, secure)(NewHandlerDeleteLesson(s)))
	http.Handle(UpsertExercise, Auth(v, rfr, secure)(NewHandlerUpsertExercise(s)))
	http.Handle(StoreExercises, Auth(v, rfr, secure)(NewHandlerStoreExercises(s)))
	http.Handle(DeleteExercise, Auth(v, rfr, secure)(NewHandlerDeleteExercise(s)))
	http.Handle(StoreResult, Auth(v, rfr, secure)(NewHandlerStoreResult(s)))

	// download
	http.Handle(ExportLessonCsv, Auth(v, rfr, secure)(NewHandlerExportLessonCsv(s)))

	// auth
	http.Handle(AuthRegister, NewHandlerAuthRegister(s, secure))
	http.Handle(AuthSignIn, HandlerNewAuthSignIn(s, secure))
	http.Handle(AuthLogout, NewHandlerAuthLogout(s))

	handler := CSRFDynamicHost()(http.DefaultServeMux)

	return http.ListenAndServe(port, handler)
}
