package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
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
			app.Button().Text("Edit").OnClick(h.onEdit(h.exercise.Id), fmt.Sprintf("%p", h)),
			app.Button().Text("Delete").OnClick(h.onDelete(h.exercise.Id), fmt.Sprintf("%p", h)),
		),
	)
}

func (h *ExerciseRow) onDelete(id int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		// delete exercise via API
		err := h.parent.api.DeleteExercise(models.Exercise{Id: id})
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
	}
}

func (h *ExerciseRow) onEdit(id int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		h.parent.inputId = h.exercise.Id
		h.parent.inputQuestion = h.exercise.Question
		h.parent.inputAnswer = h.exercise.Answer
	}
}
