package component

import (
	"fmt"
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/component/auth"
)

const PathAuthLogout = "/logout"

// Logout is a component that logs out the user
type Logout struct {
	app.Compo
	c APIClient
}

// NewLogout creates a new Logout component
func NewLogout(c APIClient) *Logout {
	return &Logout{c: c}
}

func (compo *Logout) OnMount(ctx app.Context) {
	// delete user from the local storage, so the user is logged out on the frontend
	auth.DelUser(ctx)
	// make a logout call to the backend to delete the session cookie and revoke the token
	err := compo.c.AuthLogout(ctx)
	if err != nil {
		app.Log(fmt.Errorf("failed to make a logout call: %w", err))
	}
	ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
}
