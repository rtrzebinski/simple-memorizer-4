package components

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
)

const PathAuthRegister = "/register"

// Register is a component that displays the registration form
type Register struct {
	app.Compo
	c APIClient

	// form
	inputName              string
	inputEmail             string
	inputPassword          string
	registerButtonDisabled bool
}

// NewRegister creates a new Register component
func NewRegister(c APIClient) *Register {
	return &Register{c: c}
}

// The Render method is where the component appearance is defined.
func (compo *Register) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.P().Body(
			app.Div().Body(
				app.Input().Type("text").Placeholder("Name").OnInput(compo.ValueTo(&compo.inputName)).Value(compo.inputName).Size(30),
				app.Br(),
				app.Br(),
				app.Input().Type("text").Placeholder("Email").OnInput(compo.ValueTo(&compo.inputEmail)).Value(compo.inputEmail).Size(30),
				app.Br(),
				app.Br(),
				app.Input().Type("password").Placeholder("Password").OnInput(compo.ValueTo(&compo.inputPassword)).Value(compo.inputPassword).Size(30),
				app.Br(),
				app.Br(),
				app.Button().Text("Register").OnClick(compo.handleRegister).Disabled(compo.registerButtonDisabled),
			),
		),
	)
}

// handleRegister handles register button click
func (compo *Register) handleRegister(ctx app.Context, e app.Event) {
	compo.registerButtonDisabled = true
	fmt.Println(compo.inputName)
	fmt.Println(compo.inputEmail)
	fmt.Println(compo.inputPassword)
}
