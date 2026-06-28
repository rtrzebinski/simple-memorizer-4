package http

import (
	"net/http"
	"time"
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

func NewServer(s Service, v TokenVerifier, rfr TokenRefresher, port string, secure bool) *http.Server {
	// read
	http.Handle(FetchLessons, auth(v, rfr, secure)(NewHandlerFetchLessons(s)))
	http.Handle(HydrateLesson, auth(v, rfr, secure)(NewHandlerHydrateLesson(s)))
	http.Handle(FetchExercises, auth(v, rfr, secure)(NewHandlerFetchExercises(s)))
	http.Handle(UserProfile, auth(v, rfr, secure)(NewHandlerUserProfile(v)))

	// write
	http.Handle(UpsertLesson, auth(v, rfr, secure)(NewHandlerUpsertLesson(s)))
	http.Handle(DeleteLesson, auth(v, rfr, secure)(NewHandlerDeleteLesson(s)))
	http.Handle(UpsertExercise, auth(v, rfr, secure)(NewHandlerUpsertExercise(s)))
	http.Handle(StoreExercises, auth(v, rfr, secure)(NewHandlerStoreExercises(s)))
	http.Handle(DeleteExercise, auth(v, rfr, secure)(NewHandlerDeleteExercise(s)))
	http.Handle(StoreResult, auth(v, rfr, secure)(NewHandlerStoreResult(s)))

	// download
	http.Handle(ExportLessonCsv, auth(v, rfr, secure)(NewHandlerExportLessonCsv(s)))

	// auth
	http.Handle(AuthRegister, rateLimit()(NewHandlerAuthRegister(s, secure)))
	http.Handle(AuthSignIn, rateLimit()(NewHandlerAuthSignIn(s, secure)))
	http.Handle(AuthLogout, rateLimit()(NewHandlerAuthLogout(s)))

	handler := csrf()(http.DefaultServeMux)

	return &http.Server{
		Addr:         port,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 20 * time.Second,
		IdleTimeout:  120 * time.Second,
	}
}
