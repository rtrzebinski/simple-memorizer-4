package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"github.com/rtrzebinski/simple-memorizer-4/internal/frontend"
	"net/http"
)

type Lessons struct {
	app.Compo
	api *frontend.ApiClient
}

// The OnMount method is run once component is mounted
func (c *Lessons) OnMount(ctx app.Context) {
	url := app.Window().URL()
	c.api = frontend.NewApiClient(&http.Client{}, url.Host, url.Scheme)
}

// The Render method is where the component appearance is defined.
func (c *Lessons) Render() app.UI {
	return app.Div().Body(
		&Navigation{})
}
