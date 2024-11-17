package backend

import "context"

type Service struct {
	r Reader
	w Writer
	p Publisher
}

func NewService(r Reader, w Writer, p Publisher) *Service {
	return &Service{r: r, w: w, p: p}
}

func (s *Service) FetchLessons() (Lessons, error) {
	return s.r.FetchLessons()
}

func (s *Service) HydrateLesson(lesson *Lesson) error {
	return s.r.HydrateLesson(lesson)
}

func (s *Service) FetchExercises(lesson Lesson) (Exercises, error) {
	return s.r.FetchExercises(lesson)
}

func (s *Service) UpsertLesson(lesson *Lesson) error {
	return s.w.UpsertLesson(lesson)
}

func (s *Service) DeleteLesson(lesson Lesson) error {
	return s.w.DeleteLesson(lesson)
}

func (s *Service) UpsertExercise(exercise *Exercise) error {
	return s.w.UpsertExercise(exercise)
}

func (s *Service) StoreExercises(exercises Exercises) error {
	return s.w.StoreExercises(exercises)
}

func (s *Service) DeleteExercise(exercise Exercise) error {
	return s.w.DeleteExercise(exercise)
}

func (s *Service) PublishGoodAnswer(ctx context.Context, exerciseID int) error {
	return s.p.PublishGoodAnswer(ctx, exerciseID)
}

func (s *Service) PublishBadAnswer(ctx context.Context, exerciseID int) error {
	return s.p.PublishBadAnswer(ctx, exerciseID)
}
