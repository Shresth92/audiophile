package services

import (
	"github.com/Shresth92/audiophile/services/admin"
	"github.com/Shresth92/audiophile/services/public"
	"github.com/Shresth92/audiophile/services/user"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(
		fx.Annotate(
			admin.NewAdminService,
			fx.As(
				new(AdminServices),
			),
		),
	),
	fx.Provide(
		fx.Annotate(
			public.NewPublicService,
			fx.As(
				new(PublicService),
			),
		),
	),
	fx.Provide(
		fx.Annotate(
			user.NewUserService,
			fx.As(
				new(UserServices),
			),
		),
	),
)
