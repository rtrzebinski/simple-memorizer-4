package components

import (
	"github.com/maxence-charriere/go-app/v9/pkg/app"
	"log"
)

type Home struct {
	app.Compo
	foo string
}

func NewHome(foo string) *Home {
	return &Home{
		foo: foo,
	}
}

func (h *Home) OnMount(ctx app.Context) {
	log.Println("foo: " + h.foo)
}

func (h *Home) Render() app.UI {
	return app.Div().Body()
}
