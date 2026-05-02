package auth

import (
	"context"

	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/backend"
)

type ctxKey int

const userIDKey ctxKey = 1

func ContextWithUser(ctx context.Context, user *backend.User) context.Context {
	return context.WithValue(ctx, userIDKey, user)
}

func UserIDFromContext(ctx context.Context) (userId string, ok bool) {
	user, ok := ctx.Value(userIDKey).(*backend.User)

	if !ok || user == nil {
		return "", false
	}

	return user.ID, ok
}

func UserProfileFromContext(ctx context.Context) (userProfile *backend.UserProfile, ok bool) {
	user, ok := ctx.Value(userIDKey).(*backend.User)

	if !ok || user == nil {
		return nil, false
	}

	return &backend.UserProfile{
		Name:  user.Name,
		Email: user.Email,
	}, ok
}
