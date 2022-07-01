package fx

import (
	"context"
	"github.com/qianwj/typed/collection"
	"go.uber.org/fx"
)

type Component interface {
	Provide() fx.Option
}

type Server interface {
	Component
	RegisterService(constructors ...any)
}

type DataAccess interface {
	Component
	Name(name string) DataAccess
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

type components struct {
	items []Component
}

func Components(c ...Component) Component {
	return &components{items: c}
}

func (c *components) Provide() fx.Option {
	opts := collection.Map[Component, fx.Option](c.items, func(c Component) fx.Option {
		return c.Provide()
	})
	return fx.Options(opts...)
}
