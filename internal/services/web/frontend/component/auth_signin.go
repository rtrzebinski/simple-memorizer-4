package component

import (
	"fmt"
	"log/slog"
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/component/auth"
)

const PathAuthSignIn = "/sign-in"

// SignIn is a component that displays the sign-in form
type SignIn struct {
	app.Compo
	c APIClient

	// form
	inputEmail           string
	inputPassword        string
	submitButtonDisabled bool
	errorVisible         bool
}

// NewSignIn creates a new sign-in component
func NewSignIn(c APIClient) *SignIn {
	compo := &SignIn{c: c}
	compo.inputEmail = "test.user.seeded@example.com"
	compo.inputPassword = "password"

	return compo
}

// The Render method is where the component appearance is defined.
func (compo *SignIn) Render() app.UI {
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
				app.H3().Text("Sign In"),
				app.P().Body(
					app.Text("Unable to sign in!"),
					app.Br(),
				).Style("color", "red").Hidden(!compo.errorVisible),
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

	req := frontend.SignInRequest{
		Email:    compo.inputEmail,
		Password: compo.inputPassword,
	}

	err := compo.c.AuthSignIn(ctx, req)
	if err != nil {
		compo.submitButtonDisabled = false
		compo.errorVisible = true
		app.Log(fmt.Errorf("failed to sign in: %w", err))
		return
	}

	compo.submitButtonDisabled = false
	compo.errorVisible = false
	compo.inputEmail = ""
	compo.inputPassword = ""

	user, err := compo.c.UserProfile(ctx)
	if err != nil {
		slog.Error("failed to fetch user profile", "err", err)
	}

	// persist user in local storage, so it can be used in other components
	auth.PersistUser(ctx, user)

	ctx.NavigateTo(&url.URL{Path: PathHome})
}
