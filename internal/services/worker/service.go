package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/guregu/null/v5"
)

type Service struct {
	r Reader
	w Writer
}

func NewService(r Reader, w Writer) *Service {
	return &Service{
		r: r,
		w: w,
	}
}

func (s *Service) ProcessGoodAnswer(ctx context.Context, userID string, exerciseID int) error {
	result := Result{
		Type:       Good,
		ExerciseId: exerciseID,
	}

	err := s.processAnswer(ctx, userID, result)
	if err != nil {
		return fmt.Errorf("process good answer: %w", err)
	}

	return nil
}

func (s *Service) ProcessBadAnswer(ctx context.Context, userID string, exerciseID int) error {
	result := Result{
		Type:       Bad,
		ExerciseId: exerciseID,
	}

	err := s.processAnswer(ctx, userID, result)
	if err != nil {
		return fmt.Errorf("process bad answer: %w", err)
	}

	return nil
}

func (s *Service) processAnswer(ctx context.Context, userID string, result Result) error {
	err := s.w.StoreResult(ctx, userID, result)
	if err != nil {
		return fmt.Errorf("store result: %w", err)
	}

	results, err := s.r.FetchResults(ctx, result.ExerciseId)
	if err != nil {
		return fmt.Errorf("fetch results: %w", err)
	}

	rp := resultsProjection(results)

	err = s.w.UpdateExerciseProjection(ctx, result.ExerciseId, rp)
	if err != nil {
		return fmt.Errorf("update exercise resultsProjection: %w", err)
	}

	return nil
}

func resultsProjection(results []Result) ResultsProjection {
	rp := ResultsProjection{}

	for _, a := range results {
		switch a.Type {
		case Good:
			rp.GoodAnswers++
			if isToday(a.CreatedAt) {
				rp.GoodAnswersToday++
				rp.LatestGoodAnswerWasToday = true
			}
			if a.CreatedAt.After(rp.LatestGoodAnswer.Time) {
				rp.LatestGoodAnswer = null.TimeFrom(a.CreatedAt)
			}
		case Bad:
			rp.BadAnswers++
			if isToday(a.CreatedAt) {
				rp.BadAnswersToday++
				rp.LatestBadAnswerWasToday = true
			}
			if a.CreatedAt.After(rp.LatestBadAnswer.Time) {
				rp.LatestBadAnswer = null.TimeFrom(a.CreatedAt)
			}
		}
	}

	return rp
}

func isToday(t time.Time) bool {
	now := time.Now()
	return t.Year() == now.Year() && t.Month() == now.Month() && t.Day() == now.Day()
}
