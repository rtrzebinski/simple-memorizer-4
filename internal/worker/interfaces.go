package worker

import "context"

type Reader interface {
	FetchResults(ctx context.Context, exerciseID int) ([]Result, error)
}

type Writer interface {
	StoreResult(ctx context.Context, result Result) error
	UpdateExerciseProjection(ctx context.Context, exerciseID int, projection ResultsProjection) error
}

type ProjectionBuilder interface {
	Projection(results []Result) ResultsProjection
}
