package worker

type Reader interface {
	FetchResults(exerciseID int) ([]Result, error)
}

type Writer interface {
	StoreResult(*Result) error
	UpdateExerciseProjection(exerciseID int, projection ResultsProjection) error
}

type ProjectionBuilder interface {
	Projection(results []Result) ResultsProjection
}
