package components

import (
	"fmt"
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
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
	ctx.DelState("user")
	err := compo.c.AuthLogout(ctx)
	if err != nil {
		app.Log(fmt.Errorf("failed to make a logout call: %w", err))
	}
	ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
}
