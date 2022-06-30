package fx

import (
	"context"
	"go.uber.org/fx"
)

type Component interface {
	Provide(name ...string) fx.Option
}

type Server interface {
	Component
	RegisterService(constructors ...any)
}

type DataAccess interface {
	Component
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}
