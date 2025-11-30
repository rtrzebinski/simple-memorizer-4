package components

import (
	"fmt"
	"log/slog"
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
	inputFirstName         string
	inputLastName          string
	inputEmail             string
	inputPassword          string
	registerButtonDisabled bool
}

// NewRegister creates a new Register component
func NewRegister(c APIClient) *Register {
	compo := &Register{c: c}
	compo.inputEmail = "test.user.registered@example.com"
	compo.inputPassword = "password"
	compo.inputFirstName = "Test UserProfile"
	compo.inputLastName = "Registered"

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
				app.Input().Type("text").Placeholder("First Name").OnInput(compo.ValueTo(&compo.inputFirstName)).Value(compo.inputFirstName).Size(30),
				app.Br(),
				app.Br(),
				app.Input().Type("text").Placeholder("Last Name").OnInput(compo.ValueTo(&compo.inputLastName)).Value(compo.inputLastName).Size(30),
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
	fmt.Println(compo.inputFirstName)
	fmt.Println(compo.inputLastName)
	fmt.Println(compo.inputEmail)
	fmt.Println(compo.inputPassword)

	req := frontend.RegisterRequest{
		FirstName: compo.inputFirstName,
		LastName:  compo.inputLastName,
		Email:     compo.inputEmail,
		Password:  compo.inputPassword,
	}

	err := compo.c.AuthRegister(ctx, req)
	if err != nil {
		compo.registerButtonDisabled = false
		fmt.Println("Error:", err)
		return
	}

	compo.registerButtonDisabled = false
	compo.inputFirstName = ""
	compo.inputLastName = ""
	compo.inputEmail = ""
	compo.inputPassword = ""

	user, err := compo.c.UserProfile(ctx)
	if err != nil {
		slog.Error("failed to fetch user profile", "err", err)
	}

	ctx.SetState("user", user).Persist()

	ctx.NavigateTo(&url.URL{Path: PathHome})
}
