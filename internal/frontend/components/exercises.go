package components

import (
	"fmt"
	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/backend/routes"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/api"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/csv"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/models"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend/validation"
	"net/url"
	"strconv"
)

const PathExercises = "/exercises"

// Exercises is a component that displays all exercises of a lesson
type Exercises struct {
	app.Compo
	s *api.Client

	// component vars
	lesson models.Lesson
	rows   []*ExerciseRow

	// upsert exercise form
	formVisible        bool
	validationErrors   []error
	inputId            int
	inputQuestion      string
	inputAnswer        string
	saveButtonDisabled bool
}

// NewExercises creates a new Exercises component
func NewExercises(s *api.Client) *Exercises {
	return &Exercises{
		s: s,
	}
}

// The OnMount method is run once component is mounted
func (c *Exercises) OnMount(_ app.Context) {
	u := app.Window().URL()

	// find lesson_id in the url
	lessonId, err := strconv.Atoi(u.Query().Get("lesson_id"))
	if err != nil {
		app.Log(fmt.Errorf("failed to convert lesson_id: %w", err))
		return
	}

	c.lesson = models.Lesson{Id: lessonId}

	c.hydrateLesson()
	c.displayExercisesOfLesson()
}

// hydrateLesson fetch lesson details from the database
func (c *Exercises) hydrateLesson() {
	err := c.s.HydrateLesson(&c.lesson)
	if err != nil {
		app.Log(fmt.Errorf("failed to hydrate lesson: %w", err))
	}
}

// The Render method is where the component appearance is defined.
func (c *Exercises) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.P().Body(
			app.Button().Text("Start learning").OnClick(c.handleStartLearning).Disabled(c.lesson.ExerciseCount < 2),
			app.Button().Text("Add a new exercise").OnClick(c.handleAddExercise).Hidden(c.formVisible),
			app.A().Href(routes.ExportLessonCsv+"?lesson_id="+strconv.Itoa(c.lesson.Id)).Download("").Body(
				app.Button().Text("CSV export"),
			),
			//app.Label().For("csv-upload-button").Text("CSV import"), todo style with bootstrap
			app.Input().ID("csv-upload-button").Type("file").Name("csv-import.csv").Accept(".csv").OnInput(c.handleCsvUpload),
		),
		app.P().Body(
			app.Text("Lesson name: "),
			app.Text(c.lesson.Name),
		),
		app.P().Body(
			app.Text("Lesson description: "),
			app.Text(c.lesson.Description),
		),
		app.P().Body(
			app.Text("Exercises: "),
			app.Text(c.lesson.ExerciseCount),
		),
		app.Div().Body(
			app.H3().Text("Add a new exercise"),
			app.P().Body(
				app.Range(c.validationErrors).Slice(func(i int) app.UI {
					return app.Div().Body(
						app.Text(c.validationErrors[i].Error()),
						app.Br(),
					).Style("color", "red")
				}),
			),
			app.Textarea().ID("input_question").Cols(30).Rows(3).Placeholder("Question").
				Required(true).OnInput(c.ValueTo(&c.inputQuestion)).Text(c.inputQuestion),
			app.Br(),
			app.Textarea().ID("input_answer").Cols(30).Rows(3).Placeholder("Answer").
				Required(true).OnInput(c.ValueTo(&c.inputAnswer)).Text(c.inputAnswer),
			app.Br(),
			app.Button().Text("Save").OnClick(c.handleSave).Disabled(c.saveButtonDisabled),
			app.Button().Text("Cancel").OnClick(c.handleCancel),
		).Hidden(!c.formVisible),
		app.Div().Body(
			app.H3().Text("Exercises"),
			app.Table().Style("border", "1px solid black").Body(
				&ExerciseHeader{},
				app.Range(c.rows).Slice(func(i int) app.UI {
					if c.rows[i] != nil {
						return c.rows[i].Render()
					}
					return nil
				}),
			),
		),
	)
}

// handleStartLearning start learning a current lesson
func (c *Exercises) handleStartLearning(ctx app.Context, _ app.Event) {
	u, _ := url.Parse(PathLearn)

	// set lesson_id in the url
	params := u.Query()
	params.Add("lesson_id", strconv.Itoa(c.lesson.Id))
	u.RawQuery = params.Encode()

	ctx.NavigateTo(u)
}

// handleAddExercise display add exercise form
func (c *Exercises) handleAddExercise(_ app.Context, _ app.Event) {
	c.formVisible = true
}

// handleCsvUpload upload exercises from a CSV file
func (c *Exercises) handleCsvUpload(_ app.Context, e app.Event) {
	file := e.Get("target").Get("files").Index(0)

	// validate file type which can be changed by the user in the file picker
	if file.Get("type").String() != "text/csv" {
		app.Log("invalid file type", file.Get("type").String())
		return
	}

	// read bytes from uploaded file
	data, err := readFile(file)
	if err != nil {
		app.Log(err)
		return
	}

	// extract records from bytes slice
	records, err := csv.ReadAll(data)
	if err != nil {
		app.Log(err)
		return
	}

	// prepare exercises to be stored
	var exercises models.Exercises
	for _, rec := range records {
		exercises = append(exercises, models.Exercise{
			Lesson:   &c.lesson,
			Question: rec[0],
			Answer:   rec[1],
		})
	}

	// store uploaded exercises
	err = c.s.StoreExercises(exercises)
	if err != nil {
		app.Log(err)
		return
	}

	// reset file input for next upload
	e.Get("target").Set("value", "")

	// reload the UI
	c.hydrateLesson()
	c.displayExercisesOfLesson()
}

// readFile some JS magic converting uploaded file to a slice of bytes
// https://github.com/maxence-charriere/go-app/issues/882
func readFile(file app.Value) (data []byte, err error) {
	done := make(chan bool)

	// https://developer.mozilla.org/en-US/docs/Web/API/FileReader
	reader := app.Window().Get("FileReader").New()
	reader.Set("onloadend", app.FuncOf(func(this app.Value, args []app.Value) interface{} {
		done <- true
		return nil
	}))
	reader.Call("readAsArrayBuffer", file)
	<-done

	readerError := reader.Get("error")

	if !readerError.IsNull() {
		err = fmt.Errorf("file reader error : %s", readerError.Get("message").String())
	} else {
		uint8Array := app.Window().Get("Uint8Array").New(reader.Get("result"))
		data = make([]byte, uint8Array.Length())
		app.CopyBytesToGo(data, uint8Array)
	}

	return data, err
}

// handleSave create new or update existing exercise
func (c *Exercises) handleSave(_ app.Context, e app.Event) {
	e.PreventDefault()

	// exercise to be saved
	exercise := models.Exercise{
		Id:       c.inputId,
		Question: c.inputQuestion,
		Answer:   c.inputAnswer,
		Lesson: &models.Lesson{
			Id: c.lesson.Id,
		},
	}

	// extract other questions to validate against
	var questions []string
	for i, row := range c.rows {
		if row != nil && c.rows[i].exercise.Id != c.inputId {
			questions = append(questions, c.rows[i].exercise.Question)
		}
	}

	// validate input
	validator := validation.ValidateUpsertExercise(exercise, questions)
	if validator.Failed() {
		c.validationErrors = validator.Errors()

		return
	}

	// disable submit button to avoid duplicated requests
	c.saveButtonDisabled = true

	// save exercise
	err := c.s.UpsertExercise(&exercise)
	if err != nil {
		app.Log(fmt.Errorf("failed to save exercise: %w", err))
	}

	// hide the input form on row edit, but keep open on adding new
	// because it is common to add a few exercises one after another
	if c.inputId != 0 {
		c.formVisible = false
	}

	// reset form
	c.resetForm()

	// reload the UI
	c.hydrateLesson()
	c.displayExercisesOfLesson()
}

// handleCancel handle cancel button click
func (c *Exercises) handleCancel(_ app.Context, _ app.Event) {
	c.resetForm()
	c.formVisible = false
}

// resetForm reset form fields
func (c *Exercises) resetForm() {
	c.inputId = 0
	c.inputQuestion = ""
	c.inputAnswer = ""
	c.validationErrors = nil
	c.saveButtonDisabled = false
}

// displayExercisesOfLesson fetch exercises from the database and display them
func (c *Exercises) displayExercisesOfLesson() {
	exercises, err := c.s.FetchExercises(c.lesson)
	if err != nil {
		app.Log(fmt.Errorf("failed to fetch exercises of lesson: %w", err))
	}

	// no entries in the database
	if len(exercises) == 0 {
		return
	}

	// find maxId so we know the rows slice capacity
	maxId := exercises[0].Id
	for _, row := range exercises {
		if row.Id > maxId {
			maxId = row.Id
		}
	}

	// add +1 to len as IDs from the DB are 1 indexed, while slices are 0 indexed,
	// so we need to shift by one to have space for the latest row
	c.rows = make([]*ExerciseRow, maxId+1)

	for _, row := range exercises {
		c.rows[row.Id] = &ExerciseRow{exercise: row, parent: c}
	}
}
