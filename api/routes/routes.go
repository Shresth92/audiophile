package routes

import (
	"github.com/Shresth92/audiophile/api/routes/admin"
	"github.com/Shresth92/audiophile/api/routes/public"
	"github.com/Shresth92/audiophile/api/routes/user"
	"go.uber.org/fx"
)

var Module = fx.Options(
	fx.Provide(admin.NewRoutes),
	fx.Provide(user.NewRoutes),
	fx.Provide(public.NewRoutes),
	fx.Provide(NewRoutes),
)

// Routes contains multiple routes
type Routes []Route

// Route interface
type Route interface {
	Setup()
}

// NewRoutes sets up routes
func NewRoutes(
	publicRoutes *public.Routes,
	adminRoutes *admin.Routes,
	userRoutes *user.Routes,
) *Routes {
	return &Routes{
		publicRoutes,
		adminRoutes,
		userRoutes,
	}
}

// Setup all the route
func (r *Routes) Setup() {
	for _, route := range *r {
		route.Setup()
	}
}
