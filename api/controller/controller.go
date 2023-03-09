package controller

import (
	"github.com/Shresth92/audiophile/api/controller/admin"
	"github.com/Shresth92/audiophile/api/controller/public"
	"github.com/Shresth92/audiophile/api/controller/user"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(admin.NewController),
	fx.Provide(user.NewController),
	fx.Provide(public.NewController),
)
