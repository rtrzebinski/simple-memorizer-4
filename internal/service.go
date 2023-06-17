package internal

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/projections"
	"sync"
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

	wg := sync.WaitGroup{}

	for i := range exercises {
		wg.Add(1)

		i := i

		go func() {
			answers, err := s.r.FetchAnswersOfExercise(exercises[i])
			if err != nil {
				app.Log(err)
			}

			exercises[i].AnswersProjection = projections.BuildAnswersProjection(answers)

			wg.Done()
		}()
	}

	wg.Wait()

	return exercises, nil
}

func (s *Service) FetchAnswersOfExercise(exercise models.Exercise) (models.Answers, error) {
	return s.r.FetchAnswersOfExercise(exercise)
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

func (s *Service) StoreAnswer(answer *models.Answer) error {
	return s.w.StoreAnswer(answer)
}
