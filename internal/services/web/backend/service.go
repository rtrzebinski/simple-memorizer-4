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

func (s *Service) FetchLessons(ctx context.Context) (Lessons, error) {
	return s.r.FetchLessons(ctx)
}

func (s *Service) HydrateLesson(ctx context.Context, lesson *Lesson) error {
	return s.r.HydrateLesson(ctx, lesson)
}

func (s *Service) FetchExercises(ctx context.Context, lesson Lesson) (Exercises, error) {
	return s.r.FetchExercises(ctx, lesson)
}

func (s *Service) UpsertLesson(ctx context.Context, lesson *Lesson) error {
	return s.w.UpsertLesson(ctx, lesson)
}

func (s *Service) DeleteLesson(ctx context.Context, lesson Lesson) error {
	return s.w.DeleteLesson(ctx, lesson)
}

func (s *Service) UpsertExercise(ctx context.Context, exercise *Exercise) error {
	return s.w.UpsertExercise(ctx, exercise)
}

func (s *Service) StoreExercises(ctx context.Context, exercises Exercises) error {
	return s.w.StoreExercises(ctx, exercises)
}

func (s *Service) DeleteExercise(ctx context.Context, exercise Exercise) error {
	return s.w.DeleteExercise(ctx, exercise)
}

func (s *Service) PublishGoodAnswer(ctx context.Context, exerciseID int) error {
	return s.p.PublishGoodAnswer(ctx, exerciseID)
}

func (s *Service) PublishBadAnswer(ctx context.Context, exerciseID int) error {
	return s.p.PublishBadAnswer(ctx, exerciseID)
}
