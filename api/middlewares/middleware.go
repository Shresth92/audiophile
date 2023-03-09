package middlewares

import "go.uber.org/fx"

var Module = fx.Options(
	fx.Provide(NewAdminMiddleware),
	fx.Provide(NewUserMiddleware),
	fx.Provide(NewAuthMiddleware),
)
