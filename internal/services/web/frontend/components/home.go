package components

import (
	"log/slog"
	"net/url"

	"github.com/maxence-charriere/go-app/v10/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend"
	"github.com/rtrzebinski/simple-memorizer-4/internal/services/web/frontend/components/auth"
)

const PathHome = "/"

// Home is a component that displays the home page
type Home struct {
	app.Compo
	nav *Navigation

	// component vars
	user *frontend.User
}

// NewHome creates a new Home component
func NewHome(nav *Navigation) *Home {
	return &Home{
		nav:  nav,
		user: &frontend.User{},
	}
}

// The Render method is where the component appearance is defined.
func (compo *Home) Render() app.UI {
	return app.Div().Body(
		&Navigation{},
		app.Text("Welcome "+compo.user.Name),
		app.Br(),
		app.P().Body(
			app.Text("Home page"),
		),
	)
}

// The OnMount method is run once component is mounted
func (compo *Home) OnMount(ctx app.Context) {
	// auth check
	user, err := auth.User(ctx)
	if err != nil {
		slog.Error("failed to get user", "err", err)
		ctx.NavigateTo(&url.URL{Path: PathAuthSignIn})
	}
	compo.user = user
}
