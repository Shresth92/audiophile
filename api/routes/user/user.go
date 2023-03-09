package user

import (
	"github.com/Shresth92/audiophile/api/controller/user"
	"github.com/Shresth92/audiophile/api/middlewares"
	"github.com/Shresth92/audiophile/internal"
)

type Routes struct {
	handler        *internal.RequestHandler
	controller     *user.Controller
	authMiddleware *middlewares.AuthMiddleware
	userMiddleware *middlewares.UserMiddleware
}

func NewRoutes(
	handler *internal.RequestHandler,
	controller *user.Controller,
	authMiddleware *middlewares.AuthMiddleware,
	userMiddleware *middlewares.UserMiddleware) *Routes {
	return &Routes{
		handler:        handler,
		controller:     controller,
		authMiddleware: authMiddleware,
		userMiddleware: userMiddleware,
	}
}

func (r *Routes) Setup() {
	api := r.handler.Gin.Group("/user")
	api.Use(r.authMiddleware.Setup)
	api.Use(r.userMiddleware.Setup)
	api.POST("/address", r.controller.AddAddress)
	api.GET("/offers", r.controller.GetAllOffers)

	product := api.Group("/products")
	{
		product.GET("/", r.controller.GetAllProducts)
		product.GET("/:productId", r.controller.GetProduct)

	}

	cart := api.Group("/user")
	{
		cart.GET("/", r.controller.GetMyCart)
		cart.POST("/{variantId}", r.controller.AddProductToCart)
		cart.PUT("/product-count/{variantId}", r.controller.UpdateProductCountInCart)
		cart.DELETE("/{variantId}", r.controller.RemoveProductFromCart)
		cart.DELETE("/", r.controller.DeleteMyCart)
	}

	order := api.Group("/order")
	{
		order.POST("/", r.controller.OrderProductByCart)
		order.GET("/", r.controller.GetMyOrders)
	}
}
