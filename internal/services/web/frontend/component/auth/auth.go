package auth

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
)

func PersistUser(ctx app.Context, user *frontend.UserProfile) {
	ctx.SetState("user", user).Persist()
}

func GetUser(ctx app.Context) *frontend.UserProfile {
	user := new(frontend.UserProfile)
	ctx.GetState("user", user)

	if user.Name == "" || user.Email == "" {
		return nil
	}

	return user
}

func DelUser(ctx app.Context) {
	ctx.DelState("user")
}
