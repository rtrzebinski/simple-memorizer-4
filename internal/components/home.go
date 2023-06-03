package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

var pathHome = "/"

// A Home component
type Home struct {
	app.Compo
}

// The Render method is where the component appearance is defined.
func (c *Home) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.P().Body(
			app.Text("Home page"),
		),
	)
}
