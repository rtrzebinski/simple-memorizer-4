package auth

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
)

func CheckUser(ctx app.Context) *frontend.UserProfile {
	user := new(frontend.UserProfile)
	ctx.GetState("user", user)

	if user.Name == "" || user.Email == "" {
		return nil
	}

	return user
}
