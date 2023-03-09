package internal

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewRequestHandler),
	fx.Provide(NewDatabase),
)
