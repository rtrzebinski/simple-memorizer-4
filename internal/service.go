package internal

import (
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/projections"
)

type Service struct {
	r Reader
	w Writer
}

func NewService(r Reader, w Writer) *Service {
	return &Service{r: r, w: w}
}

func (s *Service) FetchLessons() (models.Lessons, error) {
	return s.r.FetchLessons()
}

func (s *Service) HydrateLesson(lesson *models.Lesson) error {
	return s.r.HydrateLesson(lesson)
}

func (s *Service) FetchExercises(lesson models.Lesson) (models.Exercises, error) {
	exercises, err := s.r.FetchExercises(lesson)
	if err != nil {
		return exercises, err
	}

	for i := range exercises {
		exercises[i].ResultsProjection = projections.BuildResultsProjection(exercises[i].Results)
	}

	return exercises, nil
}

func (s *Service) UpsertLesson(lesson *models.Lesson) error {
	return s.w.UpsertLesson(lesson)
}

func (s *Service) DeleteLesson(lesson models.Lesson) error {
	return s.w.DeleteLesson(lesson)
}

func (s *Service) UpsertExercise(exercise *models.Exercise) error {
	return s.w.UpsertExercise(exercise)
}

func (s *Service) StoreExercises(exercises models.Exercises) error {
	return s.w.StoreExercises(exercises)
}

func (s *Service) DeleteExercise(exercise models.Exercise) error {
	return s.w.DeleteExercise(exercise)
}

func (s *Service) StoreResult(result *models.Result) error {
	return s.w.StoreResult(result)
}
