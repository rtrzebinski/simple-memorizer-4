package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type ExerciseRow struct {
	app.Compo

	exerciseId int
	question   string
	answer     string
}

// The Render method is where the component appearance is defined.
func (h *ExerciseRow) Render() app.UI {
	return app.Tr().Style("border", "1px solid black").Body(
		app.Td().Style("border", "1px solid black").Body(
			app.Text(h.question),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(h.answer),
		),
	)
}
