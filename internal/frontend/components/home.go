package components

import (
	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const PathHome = "/"

// Home is a component that displays the home page
type Home struct {
	app.Compo
}

// NewHome creates a new Home component
func NewHome() *Home {
	return &Home{}
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
