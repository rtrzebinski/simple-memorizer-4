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

func (s *Service) processAnswer(ctx context.Context, result Result) error {
	err := s.w.StoreResult(ctx, result)
	if err != nil {
		return fmt.Errorf("store result: %w", err)
	}

	results, err := s.r.FetchResults(ctx, result.ExerciseId)
	if err != nil {
		return fmt.Errorf("fetch results: %w", err)
	}

	projection := s.pb.Projection(results)

	err = s.w.UpdateExerciseProjection(ctx, result.ExerciseId, projection)
	if err != nil {
		return fmt.Errorf("update exercise projection: %w", err)
	}

	return nil
}
