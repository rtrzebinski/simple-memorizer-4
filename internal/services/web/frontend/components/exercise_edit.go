package components

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/components/auth"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/components/validation"
)

// ExerciseEdit is a component that allows creating or editing an exercise
type ExerciseEdit struct {
	app.Compo
	c      APIClient
	parent ExerciseEditor

	// form
	title              string
	visible            bool
	validationErrors   []error
	inputId            int
	inputQuestion      string
	inputAnswer        string
	saveButtonDisabled bool
}

// NewExerciseEdit creates a new ExerciseEdit component
func NewExerciseEdit(c APIClient, parent ExerciseEditor) *ExerciseEdit {
	return &ExerciseEdit{
		c:       c,
		parent:  parent,
		visible: false,
	}
}

// The Render method is where the component appearance is defined.
func (compo *ExerciseEdit) Render() app.UI {
	return app.Div().Body(
		app.H3().Text(compo.title),
		app.P().Body(
			app.Range(compo.validationErrors).Slice(func(i int) app.UI {
				return app.Div().Body(
					app.Text(compo.validationErrors[i].Error()),
					app.Br(),
				).Style("color", "red")
			}),
		),
		app.Textarea().ID("input_question").Cols(30).Rows(3).Placeholder("Question").
			Required(true).OnInput(compo.ValueTo(&compo.inputQuestion)).Text(compo.inputQuestion),
		app.Br(),
		app.Textarea().ID("input_answer").Cols(30).Rows(3).Placeholder("Answer").
			Required(true).OnInput(compo.ValueTo(&compo.inputAnswer)).Text(compo.inputAnswer),
		app.Br(),
		app.Button().Text("Save").OnClick(compo.handleSave).Disabled(compo.saveButtonDisabled),
		app.Button().Text("Cancel").OnClick(compo.handleCancel),
	).Hidden(!compo.visible)
}

// add sets the component to add mode
func (compo *ExerciseEdit) add() {
	compo.title = "Add a new exercise"
	compo.visible = true
}

// edit sets the component to edit mode with the provided exercise data
func (compo *ExerciseEdit) edit(exercise frontend.Exercise) {
	compo.title = "Edit exercise"
	compo.visible = true
	compo.inputId = exercise.Id
	compo.inputQuestion = exercise.Question
	compo.inputAnswer = exercise.Answer
	compo.validationErrors = nil
	compo.saveButtonDisabled = false
}

// handleSave create new or update existing exercise
func (compo *ExerciseEdit) handleSave(ctx app.Context, e app.Event) {
	e.PreventDefault()

	// exercise to be saved
	exercise := frontend.Exercise{
		Id:       compo.inputId,
		Question: compo.inputQuestion,
		Answer:   compo.inputAnswer,
		Lesson: &frontend.Lesson{
			Id: compo.parent.getLesson().Id,
		},
	}

	// extract other questions to validate against
	var questions []string
	for _, ex := range compo.parent.getExercises() {
		if ex != nil && ex.Id != compo.inputId {
			questions = append(questions, ex.Question)
		}
	}

	// validate input
	validator := validation.ValidateUpsertExercise(exercise, questions)
	if validator.Failed() {
		compo.validationErrors = validator.Errors()

		return
	}

	// disable submit button to avoid duplicated requests
	compo.saveButtonDisabled = true

	// save exercise
	accessToken, err := auth.Token(ctx)
	if err != nil {
		slog.Error("failed to get token", "err", err)
		ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
	}

	err = compo.c.UpsertExercise(ctx, exercise, accessToken)
	if err != nil {
		app.Log(fmt.Errorf("failed to save exercise: %w", err))
	}

	// hide the input form on row edit, but keep open on adding new
	// because it is common to add a few exercises one after another
	if compo.inputId != 0 {
		compo.visible = false
	}

	// reset form
	compo.resetForm()

	compo.parent.reloadUI(ctx)
}

// handleCancel handle cancel button click
func (compo *ExerciseEdit) handleCancel(_ app.Context, _ app.Event) {
	compo.resetForm()
	compo.visible = false
}

// resetForm reset form fields
func (compo *ExerciseEdit) resetForm() {
	compo.inputId = 0
	compo.inputQuestion = ""
	compo.inputAnswer = ""
	compo.validationErrors = nil
	compo.saveButtonDisabled = false
}
