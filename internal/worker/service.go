package worker

import (
	"context"
	"fmt"
)

type Service struct {
	r  Reader
	w  Writer
	pb ProjectionBuilder
}

func NewService(r Reader, w Writer, pb ProjectionBuilder) *Service {
	return &Service{
		r:  r,
		w:  w,
		pb: pb,
	}
}

func (s *Service) ProcessGoodAnswer(ctx context.Context, exerciseID int) error {
	result := Result{
		Type:       Good,
		ExerciseId: exerciseID,
	}

	err := s.processAnswer(ctx, result)
	if err != nil {
		return fmt.Errorf("process good answer: %w", err)
	}

	return nil
}

func (s *Service) ProcessBadAnswer(ctx context.Context, exerciseID int) error {
	result := Result{
		Type:       Bad,
		ExerciseId: exerciseID,
	}

	err := s.processAnswer(ctx, result)
	if err != nil {
		return fmt.Errorf("process bad answer: %w", err)
	}

	return nil
}

func (s *Service) processAnswer(_ context.Context, result Result) error {
	err := s.w.StoreResult(result)
	if err != nil {
		return fmt.Errorf("store result: %w", err)
	}

	results, err := s.r.FetchResults(result.ExerciseId)
	if err != nil {
		return fmt.Errorf("fetch results: %w", err)
	}

	projection := s.pb.Projection(results)

	err = s.w.UpdateExerciseProjection(result.ExerciseId, projection)
	if err != nil {
		return fmt.Errorf("update exercise projection: %w", err)
	}

	return nil
}
