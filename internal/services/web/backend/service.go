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

func (s *Service) HydrateLesson(ctx context.Context, lesson *Lesson, userID string) error {
	return s.r.HydrateLesson(ctx, lesson, userID)
}

func (s *Service) FetchExercises(ctx context.Context, lesson Lesson, oldestExerciseID int, userID string) (Exercises, error) {
	return s.r.FetchExercises(ctx, lesson, oldestExerciseID, userID)
}

func (s *Service) UpsertLesson(ctx context.Context, lesson *Lesson, userID string) error {
	return s.w.UpsertLesson(ctx, lesson, userID)
}

func (s *Service) DeleteLesson(ctx context.Context, lesson Lesson, userID string) error {
	return s.w.DeleteLesson(ctx, lesson, userID)
}

func (s *Service) UpsertExercise(ctx context.Context, exercise *Exercise, userID string) error {
	return s.w.UpsertExercise(ctx, exercise, userID)
}

func (s *Service) StoreExercises(ctx context.Context, exercises Exercises, userID string) error {
	return s.w.StoreExercises(ctx, exercises, userID)
}

func (s *Service) DeleteExercise(ctx context.Context, exercise Exercise, userID string) error {
	return s.w.DeleteExercise(ctx, exercise, userID)
}

func (s *Service) PublishGoodAnswer(ctx context.Context, exerciseID int, userID string) error {
	// todo check if userID is the owner of the exercise

	return s.p.PublishGoodAnswer(ctx, exerciseID)
}

func (s *Service) PublishBadAnswer(ctx context.Context, exerciseID int, userID string) error {
	// todo check if userID is the owner of the exercise

	return s.p.PublishBadAnswer(ctx, exerciseID)
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
