package admin

import (
	"github.com/Shresth92/audiophile/api/controller/admin"
	"github.com/Shresth92/audiophile/api/controller/user"
	"github.com/Shresth92/audiophile/api/middlewares"
	"github.com/Shresth92/audiophile/internal"
)

type Routes struct {
	handler         *internal.RequestHandler
	adminController *admin.Controller
	userController  *user.Controller
	authMiddleware  *middlewares.AuthMiddleware
	adminMiddleware *middlewares.AdminMiddleware
}

func NewRoutes(
	handler *internal.RequestHandler,
	adminController *admin.Controller,
	userController *user.Controller,
	authMiddleware *middlewares.AuthMiddleware,
	adminMiddleware *middlewares.AdminMiddleware) *Routes {
	return &Routes{
		handler:         handler,
		adminController: adminController,
		userController:  userController,
		authMiddleware:  authMiddleware,
		adminMiddleware: adminMiddleware,
	}
}

func (r *Routes) Setup() {
	api := r.handler.Gin.Group("/admin")
	api.Use(r.authMiddleware.Setup)
	api.Use(r.adminMiddleware.Setup)
	api.POST("/upload", r.adminController.UploadImages)

	products := api.Group("/products")
	{
		products.POST("/", r.adminController.CreateProduct)
		products.GET("/", r.userController.GetAllProducts)

		productId := products.Group("/:productId")
		{
			productId.GET("", r.userController.GetProduct)
			productId.PUT("", r.adminController.UpdateProduct)
			productId.DELETE("", r.adminController.DeleteProduct)

			variant := productId.Group("/variant")
			{
				variant.POST("/", r.adminController.CreateVariant)
				variant.PUT("/{variantId}", r.adminController.UpdateVariant)
				variant.DELETE("/{variantId}", r.adminController.DeleteVariant)
			}
		}
	}

	category := api.Group("/category")
	{
		category.POST("/", r.adminController.CreateCategory)
		category.GET("/", r.adminController.GetAllCategory)
		category.PUT("/{categoryId}", r.adminController.UpdateCategory)
		category.DELETE("/{categoryId}", r.adminController.DeleteCategory)
	}

	brand := api.Group("/brand")
	{
		brand.POST("/", r.adminController.CreateBrand)
		brand.GET("/", r.adminController.GetAllBrands)
		brand.PUT("/{brandId}", r.adminController.UpdateBrand)
		brand.DELETE("/{brandId}", r.adminController.DeleteBrand)
	}

	user := api.Group("/user")
	{
		user.GET("/", r.adminController.GetAllUsers)
		user.PUT("/change-role/{userId}", r.adminController.ChangeUserRole)
	}

	offer := api.Group("/offer")
	{
		offer.POST("/", r.adminController.CreateOffer)
		offer.GET("/", r.userController.GetAllOffers)
	}
}
