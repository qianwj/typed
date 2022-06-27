package grpc

import (
	"fmt"
	"go.uber.org/fx"
)

func (app *Application) RegisterRepository(constructor any, tags ...string) {
	resultTags := make([]string, len(tags))
	for i, tag := range tags {
		resultTags[i] = fmt.Sprintf(`name:"%s"`, tag)
	}
	var opt fx.Option
	if len(tags) > 0 {
		opt = fx.Provide(fx.Annotate(constructor, fx.ParamTags(resultTags...)))
	} else {
		opt = fx.Provide(constructor)
	}
	ctx := app.ctx
	ctx.repositories = append(ctx.repositories, opt)
	app.ctx = ctx
}
