package components

import (
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const PathAuthLogout = "/logout"

// Logout is a component that logs out the user
type Logout struct {
	app.Compo
}

// NewLogout creates a new Logout component
func NewLogout() *Logout {
	return &Logout{}
}

func (compo *Logout) OnMount(ctx app.Context) {
	ctx.DelState("AccessToken")
	ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
}
