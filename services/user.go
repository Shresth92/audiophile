package services

import "github.com/Shresth92/audiophile/models"

type UserServices interface {
	AddProductToCart(userID string, variantID string) error
	UpdateProductCountInCart(userID string, variantID string, count bool) error
	RemoveProductFromCart(userID string, variantID string) error
	DeleteCart(userID string) error
	GetCartProducts(userID string) ([]models.UserCart, error)
	GetProduct(productId string) ([]models.AllProducts, error)
	GetTotalProductCost(variantIds []string) (int, error)
	PriceAfterDiscount(price int, couponCode string) (int, error)
	FilterMyOrders(userId string, ProductStatus models.DeliveryStatus, limit int, page int) ([]models.Orders, error)
	CountFilterMyOrders(userId string, ProductStatus models.DeliveryStatus) (int64, error)
	UpdateProductStock(variantId string, stock int) error
	FilterAllProducts(limit int, page int, searchString string, categoryFilter string, brandFilter string) ([]models.AllProducts, error)
	FilterAllProductsCount(searchString string, categoryFilter string, brandFilter string) (int64, error)
	GenerateOrderIdByCart(price int, userId string, addressId string) (string, error)
	AddProductsInOrder(orderedProducts []models.ProductOrdered) error
	GetAllOffers() ([]models.Offer, error)
	AddAddress(userID string, newAddress *models.Address) (string, error)
}
