package fx

import (
	"context"
	"go.uber.org/fx"
)

type DataAccess interface {
	Provide(name ...string) fx.Option
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}
