package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
)

// ExerciseRow is a component that displays a single row of the exercises table
type ExerciseRow struct {
	app.Compo
	parent   *Exercises
	exercise models.Exercise
}

// The Render method is where the component appearance is defined.
func (compo *ExerciseRow) Render() app.UI {
	return app.Tr().Style("border", "1px solid black").Body(
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.exercise.Id),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.exercise.Question),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.exercise.Answer),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.exercise.ResultsProjection.BadAnswers),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.exercise.ResultsProjection.GoodAnswers),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(compo.exercise.ResultsProjection.GoodAnswersPercent()),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Button().Text("Edit").OnClick(compo.onEdit(), app.EventScope(fmt.Sprintf("%p", compo))),
			app.Button().Text("Delete").OnClick(compo.onDelete(compo.exercise.Id), app.EventScope(fmt.Sprintf("%p", compo))),
		),
	)
}

// onDelete handles delete button click
func (compo *ExerciseRow) onDelete(id int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		// delete exercise
		err := compo.parent.c.DeleteExercise(models.Exercise{Id: id})
		if err != nil {
			app.Log(fmt.Errorf("failed to delete exercise: %w", err))
		}
		// create a new rows slice to be replaced in parent component
		rows := make([]*ExerciseRow, len(compo.parent.rows))
		for i, row := range compo.parent.rows {
			// add all rows but current one (which is being deleted)
			if i != id && row != nil {
				rows[i] = row
			}
		}
		// replace parent rows slice with a new one - this will update the UI
		compo.parent.rows = rows
		compo.parent.hydrateLesson()
	}
}

// onEdit handles edit button click
func (compo *ExerciseRow) onEdit() app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		compo.parent.inputId = compo.exercise.Id
		compo.parent.inputQuestion = compo.exercise.Question
		compo.parent.inputAnswer = compo.exercise.Answer
		compo.parent.formVisible = true
		compo.parent.validationErrors = nil
	}
}
