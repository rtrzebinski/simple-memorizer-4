package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type ExerciseHeader struct {
	app.Compo
}

// The Render method is where the component appearance is defined.
func (c *ExerciseHeader) Render() app.UI {
	return app.Tr().Style("border", "1px solid black").Body(
		app.Th().Style("border", "1px solid black").Body(
			app.Text("Id"),
		),
		app.Th().Style("border", "1px solid black").Body(
			app.Text("Question"),
		),
		app.Th().Style("border", "1px solid black").Body(
			app.Text("Answer"),
		),
		app.Th().Style("border", "1px solid black").Body(
			app.Text("Bad answers"),
		),
		app.Th().Style("border", "1px solid black").Body(
			app.Text("Good answers"),
		),
		app.Th().Style("border", "1px solid black").Body(
			app.Text("Good answers %"),
		),
		app.Th().Style("border", "1px solid black").Body(
			app.Text("Actions"),
		),
	)
}
