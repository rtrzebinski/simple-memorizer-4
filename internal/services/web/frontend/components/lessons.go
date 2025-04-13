package components

import (
	"context"
	"fmt"
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/components/auth"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/components/validation"
)

const PathLessons = "/lessons"

// Lessons is a component that displays all lessons
type Lessons struct {
	app.Compo
	c    APIClient
	user *frontend.User

	// component vars
	rows []*LessonRow

	// save lesson form
	formVisible        bool
	validationErrors   []error
	inputId            int
	inputName          string
	inputDescription   string
	saveButtonDisabled bool
}

// NewLessons creates a new Lessons component
func NewLessons(c APIClient) *Lessons {
	return &Lessons{
		c:    c,
		user: &frontend.User{},
	}
}

// The OnMount method is run once component is mounted
func (compo *Lessons) OnMount(ctx app.Context) {
	// auth check
	user, err := auth.User(ctx)
	if err != nil {
		ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
	}
	compo.user = user

	compo.displayAllLessons(ctx)
}

// The Render method is where the component appearance is defined.
func (compo *Lessons) Render() app.UI {
	return app.Div().Body(
		app.Div().Body(
			app.P().Body(
				app.A().Href(PathHome).Text("Home"),
				app.Text(" | "),
				app.A().Href(PathLessons).Text("Lessons"),
				app.Text(" | "),
				app.A().Href(PathAuthLogout).Text("Logout"),
				app.Text(" | "),
				app.Text(app.Getenv("version")),
			),
		),
		app.Text("Welcome "+compo.user.Name),
		app.Br(),
		app.P().Body(
			app.Button().Text("Add a new lesson").OnClick(compo.handleAddLesson).Hidden(compo.formVisible),
		),
		app.Div().Body(
			app.H3().Text("Add a new lesson"),
			app.P().Body(
				app.Range(compo.validationErrors).Slice(func(i int) app.UI {
					return app.Div().Body(
						app.Text(compo.validationErrors[i].Error()),
						app.Br(),
					).Style("color", "red")
				}),
			),
			app.Textarea().Cols(30).Rows(3).Placeholder("Name").
				Required(true).OnInput(compo.ValueTo(&compo.inputName)).Text(compo.inputName),
			app.Br(),
			app.Textarea().Cols(30).Rows(3).Placeholder("Description").
				Required(true).OnInput(compo.ValueTo(&compo.inputDescription)).Text(compo.inputDescription),
			app.Br(),
			app.Button().Text("Save").OnClick(compo.handleSave).Disabled(compo.saveButtonDisabled),
			app.Button().Text("Cancel").OnClick(compo.handleCancel),
		).Hidden(!compo.formVisible),
		app.Div().Body(
			app.H3().Text("Lessons"),
			app.Table().Style("border", "1px solid black").Body(
				&LessonHeader{},
				app.Range(compo.rows).Slice(func(i int) app.UI {
					if compo.rows[i] != nil {
						return compo.rows[i].Render()
					}
					return nil
				}),
			),
		),
	)
}

// handleAddLesson display add lesson form
func (compo *Lessons) handleAddLesson(_ app.Context, e app.Event) {
	compo.formVisible = true
	compo.inputId = 0
}

// handleSave create new or update existing lesson
func (compo *Lessons) handleSave(ctx app.Context, e app.Event) {
	e.PreventDefault()

	var err error

	// lesson to be saved
	lesson := frontend.Lesson{
		Id:          compo.inputId,
		Name:        compo.inputName,
		Description: compo.inputDescription,
	}

	// extract other names to validate against
	var names []string
	for i, row := range compo.rows {
		if row != nil && compo.rows[i].lesson.Id != compo.inputId {
			names = append(names, compo.rows[i].lesson.Name)
		}
	}

	// validate input
	validator := validation.ValidateUpsertLesson(lesson, names)
	if validator.Failed() {
		compo.validationErrors = validator.Errors()

		return
	}

	// disable submit button to avoid duplicated requests
	compo.saveButtonDisabled = true

	// save lesson
	err = compo.c.UpsertLesson(ctx, lesson)
	if err != nil {
		app.Log(fmt.Errorf("failed to save lesson: %w", err))
	}

	// reset form
	compo.resetForm()

	// refresh lessons list
	compo.displayAllLessons(ctx)
}

// handleCancel handle cancel button click
func (compo *Lessons) handleCancel(_ app.Context, _ app.Event) {
	compo.resetForm()
}

// resetForm reset form fields
func (compo *Lessons) resetForm() {
	compo.inputId = 0
	compo.inputName = ""
	compo.inputDescription = ""
	compo.validationErrors = nil
	compo.saveButtonDisabled = false
	compo.formVisible = false
}

// displayAllLessons fetch all lessons and display them
func (compo *Lessons) displayAllLessons(ctx context.Context) {
	lessons, err := compo.c.FetchLessons(ctx)
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch all lessons: %w", err))
	}

	// no entries
	if len(lessons) == 0 {
		return
	}

	// find maxId so we know the rows slice capacity
	maxId := lessons[0].Id
	for _, row := range lessons {
		if row.Id > maxId {
			maxId = row.Id
		}
	}

	// add +1 to len as IDs from the DB are 1 indexed, while slices are 0 indexed,
	// so we need to shift by one to have space for the latest row
	compo.rows = make([]*LessonRow, maxId+1)

	for _, row := range lessons {
		compo.rows[row.Id] = &LessonRow{lesson: row, parent: compo}
	}
}
