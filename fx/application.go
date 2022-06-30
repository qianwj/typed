package fx

import (
	"github.com/qianwj/typed/fx/options"
)

type Application interface {
	RegisterService(constructors ...any)
	RegisterRepository(constructor any, tags ...string)
	WithData(data DataAccess, name ...string) Application
	WithLogger(constructor any) Application
	Run()
}

func NewApp() Application {
	return &defaultApp{}
}

type defaultApp struct {
	opts []options.Options
}

func (app *defaultApp) Run() {}

func (app *defaultApp) RegisterRepository(constructor any, tags ...string) {}

func (app *defaultApp) WithData(data DataAccess, name ...string) Application {
	return app
}

func (app *defaultApp) WithLogger(constructor any) Application {
	return app
}

func (app *defaultApp) RegisterService(constructors ...any) {}
