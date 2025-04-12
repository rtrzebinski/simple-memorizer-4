package components

import (
	"log/slog"
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const PathHome = "/"

// Home is a component that displays the home page
type Home struct {
	app.Compo
	nav *Navigation
}

// NewHome creates a new Home component
func NewHome(nav *Navigation) *Home {
	return &Home{nav: nav}
}

// The Render method is where the component appearance is defined.
func (compo *Home) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.P().Body(
			app.Text("Home page"),
		),
	)
}

// The OnMount method is run once component is mounted
func (compo *Home) OnMount(ctx app.Context) {
	// auth check
	var at string
	ctx.GetState("resp.AccessToken", &at)
	if at == "" {
		ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
	} else {
		compo.nav.showLessons = true
		compo.nav.showHome = true
		slog.Info("access token", "token", at)
	}
}
