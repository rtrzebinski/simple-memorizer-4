package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
	"strconv"
)

type ExerciseRow struct {
	app.Compo
	parent   *Exercises
	exercise models.Exercise
}

// The Render method is where the component appearance is defined.
func (h *ExerciseRow) Render() app.UI {
	return app.Tr().Style("border", "1px solid black").Body(
		app.Td().Style("border", "1px solid black").Body(
			app.Text(h.exercise.Id),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(h.exercise.Question),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(h.exercise.Answer),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(h.exercise.BadAnswers),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(h.exercise.GoodAnswers),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Button().ID(strconv.Itoa(h.exercise.Id)).Text("Delete").OnClick(func(ctx app.Context, e app.Event) {
				// get exercise id to be deleted from the button ID,
				// we can not use the "h.exercise.Id" because of the bug:
				// https://github.com/maxence-charriere/go-app/issues/826
				id, err := strconv.Atoi(ctx.JSSrc().Get("id").String())
				if err != nil {
					app.Log(fmt.Errorf("failed to convert row id to int: %w", err))
				}
				// delete exercise via API
				err = h.parent.api.DeleteExercise(models.Exercise{Id: id})
				if err != nil {
					app.Log(fmt.Errorf("failed to delete exercise: %w", err))
				}
				// create a new rows slice to be replaced in parent component
				rows := make([]*ExerciseRow, len(h.parent.rows))
				for i, row := range h.parent.rows {
					// add all rows but current one (which is being deleted)
					if i != id && row != nil {
						rows[i] = row
					}
				}
				// replace parent rows slice with a new one - this will update the UI
				h.parent.rows = rows
			}),
		),
	)
}
