package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type LessonHeader struct {
	app.Compo
}

// The Render method is where the component appearance is defined.
func (c *LessonHeader) Render() app.UI {
	return app.Tr().Style("border", "1px solid black").Body(
		app.Th().Style("border", "1px solid black").Body(
			app.Text("Id"),
		),
		app.Th().Style("border", "1px solid black").Body(
			app.Text("Name"),
		),
		app.Th().Style("border", "1px solid black").Body(
			app.Text("Description"),
		),
		app.Th().Style("border", "1px solid black").Body(
			app.Text("Exercises"),
		),
		app.Th().Style("border", "1px solid black").Body(
			app.Text("Actions"),
		),
	)
}
