package public

import (
	"github.com/Shresth92/audiophile/api/controller/public"
	"github.com/Shresth92/audiophile/internal"
)

type Routes struct {
	handler    *internal.RequestHandler
	controller *public.Controller
}

func NewRoutes(
	handler *internal.RequestHandler,
	controller *public.Controller) *Routes {
	return &Routes{
		handler:    handler,
		controller: controller,
	}
}

func (r *Routes) Setup() {
	api := r.handler.Gin.Group("/public")
	api.POST("/register", r.controller.Register)
	api.POST("/login", r.controller.Login)
	api.POST("/admin-login", r.controller.Login)
}
