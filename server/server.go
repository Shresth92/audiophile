package server

import (
	"context"
	"github.com/Shresth92/audiophile/handler"
	"github.com/Shresth92/audiophile/middlewares"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

func ServeStart() *Server {
	router := chi.NewRouter()
	router.Route("/audiophile", func(audiophile chi.Router) {
		audiophile.Route("/public", func(public chi.Router) {
			public.Post("/register", handler.Register)
			public.Post("/login", handler.Login)
			public.Post("/admin-login", handler.Login)
		})
		audiophile.Route("/admin", func(admin chi.Router) {
			admin.Use(middlewares.AuthMiddleware)
			admin.Use(middlewares.CheckAdmin)
			admin.Post("/logout", handler.Logout)
			admin.Post("/upload", handler.UploadImages)
			admin.Route("/products", func(products chi.Router) {
				products.Post("/", handler.CreateProduct)
				products.Get("/all", handler.GetAllProducts)
				products.Get("/{productId}", handler.GetProduct)
				products.Put("/{productId}", handler.UpdateProduct)
				products.Delete("/{productId}", handler.DeleteProduct)
				products.Route("/{productId}/variants", func(variant chi.Router) {
					variant.Post("/", handler.CreateVariant)
					variant.Put("/{variantId}", handler.UpdateVariant)
					variant.Delete("/{variantId}", handler.DeleteVariant)
				})

			})
			admin.Route("/category", func(category chi.Router) {
				category.Post("/", handler.CreateCategory)
				category.Get("/", handler.GetAllCategory)
				category.Put("/{categoryId}", handler.UpdateCategory)
				category.Delete("/{categoryId}", handler.DeleteCategory)
			})
			admin.Route("/brand", func(brand chi.Router) {
				brand.Post("/", handler.CreateBrand)
				brand.Get("/", handler.GetAllBrands)
				brand.Put("/{brandId}", handler.UpdateBrand)
				brand.Delete("/{brandId}", handler.DeleteBrand)
			})
			admin.Route("/user", func(user chi.Router) {
				user.Get("/", handler.GetAllUsers)
				user.Put("/change-role/{userId}", handler.ChangeUserRole)
			})
			admin.Route("/offer", func(offer chi.Router) {
				offer.Post("/", handler.CreateOffer)
				offer.Get("/", handler.GetAllOffers)
			})
		})
		audiophile.Route("/user", func(user chi.Router) {
			user.Use(middlewares.AuthMiddleware)
			user.Use(middlewares.CheckUser)
			user.Post("/logout", handler.Logout)
			user.Post("/address", handler.AddAddress)
			user.Get("/offers", handler.GetAllOffers)
			user.Route("/products", func(product chi.Router) {
				product.Get("/", handler.GetAllProducts)
				product.Get("/{productId}", handler.GetProduct)
				product.Put("/count-in-cart/{variantId}", handler.UpdateProductCountInCart)
			})
			user.Route("/cart", func(cart chi.Router) {
				cart.Get("/", handler.GetMyCart)
				cart.Post("/{variantId}", handler.AddProductToCart)
				cart.Delete("/{variantId}", handler.RemoveProductFromCart)
				cart.Delete("/", handler.DeleteMyCart)
			})
			user.Route("/order", func(order chi.Router) {
				order.Post("/", handler.OrderProductByCart)
				order.Get("/", handler.GetMyOrders)
			})
		})
	})

	return &Server{
		Router: router,
	}
}

func (s *Server) Start(port string) error {
	s.server = &http.Server{
		Addr:         port,
		Handler:      s.Router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	log.Println("Server started at ", port)
	return s.server.ListenAndServe()
}

func (s *Server) Stop(time time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), time)
	defer cancel()
	return s.server.Shutdown(ctx)
}
