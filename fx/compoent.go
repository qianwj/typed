package fx

import (
	"context"
	"github.com/qianwj/typed/collection"
	"go.uber.org/fx"
	"sync"
)

type Component interface {
	Provide() fx.Option
}

type OrderedComponent interface {
	Component
	Order() int
}

type Server interface {
	Component
	RegisterService(constructors ...any)
}

type DataSource interface {
	Component
	Name(name string) DataSource
	Connect(ctx context.Context) error
	Close(ctx context.Context) error
}

type DataSources struct {
	fx.In
	locker *sync.Mutex
	Items  []DataSource `group:"data_sources"`
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

func (d *DataSources) Connect(ctx context.Context) error {
	d.locker.Lock()
	for _, dataSource := range d.Items {
		if err := dataSource.Connect(ctx); err != nil {
			d.locker.Unlock()
			return err
		}
	}
	d.locker.Unlock()
	return nil
}

func (d *DataSources) Close(ctx context.Context) error {
	d.locker.Lock()
	var err error
	for _, dataSource := range d.Items {
		err = dataSource.Close(ctx)
	}
	d.locker.Unlock()
	return err
}
