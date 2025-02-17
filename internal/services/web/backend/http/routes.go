package http

const (
	// read

	FetchLessons   = "/fetch-lessons"
	HydrateLesson  = "/hydrate-lesson"
	FetchExercises = "/fetch-exercises"

	// write

	UpsertLesson   = "/upsert-lesson"
	DeleteLesson   = "/delete-lesson"
	UpsertExercise = "/upsert-exercise"
	StoreExercises = "/store-exercises"
	DeleteExercise = "/delete-exercise"
	StoreResult    = "/store-result"

	// download

	ExportLessonCsv = "/export-lesson-csv"
)
