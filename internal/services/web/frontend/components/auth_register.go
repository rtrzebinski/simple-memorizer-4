package components

import (
	"fmt"
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
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
	compo := &Register{c: c}
	compo.inputEmail = "foo@bar.com"
	compo.inputPassword = "password"
	compo.inputName = "foo"

	return compo
}

// The Render method is where the component appearance is defined.
func (compo *Register) Render() app.UI {
	return app.Div().Body(
		app.Div().Body(
			app.P().Body(
				app.A().Href(PathAuthSignIn).Text("SignIn"),
				app.Text(" | "),
				app.A().Href(PathAuthRegister).Text("Register"),
				app.Text(" | "),
				app.Text(app.Getenv("version")),
			),
		),
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

	req := frontend.RegisterRequest{
		Name:     compo.inputName,
		Email:    compo.inputEmail,
		Password: compo.inputPassword,
	}

	resp, err := compo.c.AuthRegister(ctx, req)
	if err != nil {
		compo.registerButtonDisabled = false
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Response:", resp)
	compo.registerButtonDisabled = false
	compo.inputName = ""
	compo.inputEmail = ""
	compo.inputPassword = ""

	ctx.SetState("AccessToken", resp.AccessToken).Persist()
	ctx.NavigateTo(&url.URL{Path: PathHome})
}
