package components

import (
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const PathAuthLogout = "/logout"

// Logout is a component that logs out the user
type Logout struct {
	app.Compo
	nav *Navigation
}

// NewLogout creates a new Logout component
func NewLogout(nav *Navigation) *Logout {
	return &Logout{nav: nav}
}

func (compo *Logout) OnMount(ctx app.Context) {
	ctx.DelState("resp.AccessToken")
	ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
	compo.nav.showLessons = false
	compo.nav.showHome = false
}
