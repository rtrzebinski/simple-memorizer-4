package backend

import "context"

type Service struct {
	r Reader
	w Writer
	p Publisher
	a AuthClient
}

func NewService(r Reader, w Writer, p Publisher, a AuthClient) *Service {
	return &Service{r: r, w: w, p: p, a: a}
}

func (s *Service) FetchLessons(ctx context.Context, userID string) (Lessons, error) {
	return s.r.FetchLessons(ctx, userID)
}

func (s *Service) HydrateLesson(ctx context.Context, userID string, lesson *Lesson) error {
	return s.r.HydrateLesson(ctx, userID, lesson)
}

func (s *Service) FetchExercises(ctx context.Context, userID string, lesson Lesson, oldestExerciseID int) (Exercises, error) {
	return s.r.FetchExercises(ctx, userID, lesson, oldestExerciseID)
}

func (s *Service) UpsertLesson(ctx context.Context, userID string, lesson *Lesson) error {
	return s.w.UpsertLesson(ctx, userID, lesson)
}

func (s *Service) DeleteLesson(ctx context.Context, userID string, lesson Lesson) error {
	return s.w.DeleteLesson(ctx, userID, lesson)
}

func (s *Service) UpsertExercise(ctx context.Context, userID string, exercise *Exercise) error {
	return s.w.UpsertExercise(ctx, userID, exercise)
}

func (s *Service) StoreExercises(ctx context.Context, userID string, exercises Exercises) error {
	return s.w.StoreExercises(ctx, userID, exercises)
}

func (s *Service) DeleteExercise(ctx context.Context, userID string, exercise Exercise) error {
	return s.w.DeleteExercise(ctx, userID, exercise)
}

func (s *Service) PublishGoodAnswer(ctx context.Context, userID string, exerciseID int) error {
	return s.p.PublishGoodAnswer(ctx, userID, exerciseID)
}

func (s *Service) PublishBadAnswer(ctx context.Context, userID string, exerciseID int) error {
	return s.p.PublishBadAnswer(ctx, userID, exerciseID)
}

func (s *Service) Register(ctx context.Context, firstName, lastName, email, password string) (Tokens, error) {
	return s.a.Register(ctx, firstName, lastName, email, password)
}

func (s *Service) SignIn(ctx context.Context, email, password string) (Tokens, error) {
	return s.a.SignIn(ctx, email, password)
}

func (s *Service) Revoke(ctx context.Context, refreshToken string) error {
	return s.a.Revoke(ctx, refreshToken)
}
