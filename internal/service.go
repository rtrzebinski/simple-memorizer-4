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

func (s *Service) FetchAllLessons() (models.Lessons, error) {
	return s.r.FetchAllLessons()
}

func (s *Service) HydrateLesson(lesson *models.Lesson) error {
	return s.r.HydrateLesson(lesson)
}

func (s *Service) FetchExercisesOfLesson(lesson models.Lesson) (models.Exercises, error) {
	exercises, err := s.r.FetchExercisesOfLesson(lesson)
	if err != nil {
		return exercises, err
	}

	for i := range exercises {
		exercises[i].ResultsProjection = projections.BuildResultsProjection(exercises[i].Results)
	}

	return exercises, nil
}

func (s *Service) FetchResultsOfExercise(exercise models.Exercise) (models.Results, error) {
	return s.r.FetchResultsOfExercise(exercise)
}

func (s *Service) StoreLesson(lesson *models.Lesson) error {
	return s.w.StoreLesson(lesson)
}

func (s *Service) DeleteLesson(lesson models.Lesson) error {
	return s.w.DeleteLesson(lesson)
}

func (s *Service) StoreExercise(exercise *models.Exercise) error {
	return s.w.StoreExercise(exercise)
}

func (s *Service) DeleteExercise(exercise models.Exercise) error {
	return s.w.DeleteExercise(exercise)
}

func (s *Service) StoreResult(result *models.Result) error {
	return s.w.StoreResult(result)
}
