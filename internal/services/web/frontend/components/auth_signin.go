package components

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const PathAuthSignIn = "/sign-in"

// SignIn is a component that displays the sign-in form
type SignIn struct {
	app.Compo
	c APIClient

	// form
	inputName            string
	inputEmail           string
	inputPassword        string
	submitButtonDisabled bool
}

// NewSignIn creates a new sign-in component
func NewSignIn(c APIClient) *SignIn {
	return &SignIn{c: c}
}

// The Render method is where the component appearance is defined.
func (compo *SignIn) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.P().Body(
			app.Div().Body(
				app.Input().Type("text").Placeholder("Email").OnInput(compo.ValueTo(&compo.inputEmail)).Value(compo.inputEmail).Size(30),
				app.Br(),
				app.Br(),
				app.Input().Type("password").Placeholder("Password").OnInput(compo.ValueTo(&compo.inputPassword)).Value(compo.inputPassword).Size(30),
				app.Br(),
				app.Br(),
				app.Button().Text("Submit").OnClick(compo.handleSubmit).Disabled(compo.submitButtonDisabled),
			),
		),
	)
}

// handleSubmit handles submit button click
func (compo *SignIn) handleSubmit(ctx app.Context, e app.Event) {
	compo.submitButtonDisabled = true
	fmt.Println(compo.inputEmail)
	fmt.Println(compo.inputPassword)
}
