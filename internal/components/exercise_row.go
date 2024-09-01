package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/models"
)

type ExerciseRow struct {
	app.Compo
	parent   *Exercises
	exercise models.Exercise
}

// The Render method is where the component appearance is defined.
func (c *ExerciseRow) Render() app.UI {
	return app.Tr().Style("border", "1px solid black").Body(
		app.Td().Style("border", "1px solid black").Body(
			app.Text(c.exercise.Id),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(c.exercise.Question),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(c.exercise.Answer),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(c.exercise.ResultsProjection.BadAnswers),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(c.exercise.ResultsProjection.GoodAnswers),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Text(c.exercise.ResultsProjection.GoodAnswersPercent()),
		),
		app.Td().Style("border", "1px solid black").Body(
			app.Button().Text("Edit").OnClick(c.onEdit()),
			app.Button().Text("Delete").OnClick(c.onDelete(c.exercise.Id)),
		),
	)
}

func (c *ExerciseRow) onDelete(id int) app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		// delete exercise
		err := c.parent.s.DeleteExercise(models.Exercise{Id: id})
		if err != nil {
			app.Log(fmt.Errorf("failed to delete exercise: %w", err))
		}
		// create a new rows slice to be replaced in parent component
		rows := make([]*ExerciseRow, len(c.parent.rows))
		for i, row := range c.parent.rows {
			// add all rows but current one (which is being deleted)
			if i != id && row != nil {
				rows[i] = row
			}
		}
		// replace parent rows slice with a new one - this will update the UI
		c.parent.rows = rows
		c.parent.hydrateLesson()
	}
}

func (c *ExerciseRow) onEdit() app.EventHandler {
	return func(ctx app.Context, e app.Event) {
		c.parent.inputId = c.exercise.Id
		c.parent.inputQuestion = c.exercise.Question
		c.parent.inputAnswer = c.exercise.Answer
		c.parent.formVisible = true
		c.parent.validationErrors = nil
	}
}
