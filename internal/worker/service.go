package worker

import (
	"context"
	"fmt"
)

type Service struct {
	w Writer
}

func NewService(w Writer) *Service {
	return &Service{
		w: w,
	}
}

func (s *Service) ProcessGoodAnswer(_ context.Context, exerciseID int) error {
	result := &Result{
		Type:       Good,
		ExerciseId: exerciseID,
	}

	err := s.w.StoreResult(result)
	if err != nil {
		return fmt.Errorf("store good result: %w", err)
	}

	fmt.Printf("good answer processed %d\n", exerciseID)

	return nil
}

func (s *Service) ProcessBadAnswer(_ context.Context, exerciseID int) error {
	result := &Result{
		Type:       Bad,
		ExerciseId: exerciseID,
	}

	err := s.w.StoreResult(result)
	if err != nil {
		return fmt.Errorf("store bad result: %w", err)
	}

	fmt.Printf("bad answer processed %d\n", exerciseID)

	return nil
}
