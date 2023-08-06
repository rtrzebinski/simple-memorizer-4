package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Navigation struct {
	app.Compo
}

// The Render method is where the component appearance is defined.
func (c *Navigation) Render() app.UI {
	return app.Div().Body(
		app.P().Body(
			app.A().Href(PathHome).Text("Home"),
			app.Text(" | "),
			app.A().Href(PathLessons).Text("Lessons"),
			app.Text(" | "),
			app.Text(app.Getenv("version")),
		),
	)
}
