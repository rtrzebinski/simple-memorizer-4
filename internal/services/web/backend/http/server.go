package http

import (
	"net/http"
	"time"

	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"github.com/ulule/limiter/v3/drivers/store/memory"
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
	http.Handle(FetchExercises, Auth(v, rfr, secure)(NewHandlerFetchExercises(s)))
	http.Handle(UserProfile, Auth(v, rfr, secure)(NewHandlerUserProfile(v)))

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
	http.Handle(AuthRegister, rateLimit(NewHandlerAuthRegister(s, secure)))
	http.Handle(AuthSignIn, rateLimit(NewHandlerAuthSignIn(s, secure)))
	http.Handle(AuthLogout, rateLimit(NewHandlerAuthLogout(s)))

	handler := CSRFDynamicHost()(http.DefaultServeMux)

	return http.ListenAndServe(port, handler)
}

// rateLimit takes an http.Handler and returns a new http.Handler wrapped with rate limiting
func rateLimit(next http.Handler) http.Handler {
	rate := limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  20,
	}

	store := memory.NewStoreWithOptions(limiter.StoreOptions{
		CleanUpInterval: 5 * time.Minute,
	})

	// K8s Ingress must be configured to set the X-Forwarded-For header
	// https://github.com/ulule/limiter#limiter-behind-a-reverse-proxy
	instance := limiter.New(store, rate, limiter.WithTrustForwardHeader(true))

	m := stdlib.NewMiddleware(instance)
	m.OnLimitReached = func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
	}

	return m.Handler(next)
}
