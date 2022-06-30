package fx

import (
	"github.com/qianwj/typed/collection"
	"go.uber.org/fx"
)

type Application struct {
	components []fx.Option
}

func NewApp(components ...Component) *Application {
	opts := collection.Map[Component, fx.Option](components, func(c Component) fx.Option {
		return c.Provide()
	})
	return &Application{components: opts}
}

func (app *Application) Provide(components ...fx.Option) *Application {
	app.components = append(app.components, components...)
	return app
}

func (app *Application) Run() {
	container := fx.New(app.components...)
	container.Run()
}
