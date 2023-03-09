package user

import (
	"github.com/Shresth92/audiophile/internal"
	"github.com/Shresth92/audiophile/models"
)

type Service struct {
	repo *repository
}

func NewUserService(db *internal.Database) *Service {
	return &Service{repo: newUserRepository(db)}
}

func (s *Service) AddProductToCart(userID string, variantID string) error {
	return s.repo.addProductToCart(userID, variantID, 1)
}

func (s *Service) UpdateProductCountInCart(userID string, variantID string, count bool) error {
	return s.repo.updateProductCountInCart(userID, variantID, count)
}

func (s *Service) RemoveProductFromCart(userID string, variantID string) error {
	return s.repo.removeCartProduct(userID, variantID)
}

func (s *Service) DeleteCart(userID string) error {
	return s.repo.deleteCart(userID)
}

func (s *Service) GetCartProducts(userID string) ([]models.UserCart, error) {
	return s.repo.getCartProducts(userID)
}

func (s *Service) GetProduct(productId string) ([]models.AllProducts, error) {
	return s.repo.getProduct(productId)
}

func (s *Service) GetTotalProductCost(variantIds []string) (int, error) {
	return s.repo.getTotalProductCost(variantIds)
}

func (s *Service) PriceAfterDiscount(price int, couponCode string) (int, error) {
	return s.repo.priceAfterDiscount(price, couponCode)
}

func (s *Service) FilterMyOrders(userId string, ProductStatus models.DeliveryStatus, limit int, page int) ([]models.Orders, error) {
	return s.repo.filterMyOrders(userId, ProductStatus, limit, page)
}

func (s *Service) CountFilterMyOrders(userId string, ProductStatus models.DeliveryStatus) (int64, error) {
	return s.repo.countFilterMyOrders(userId, ProductStatus)
}

func (s *Service) UpdateProductStock(variantId string, stock int) error {
	return s.repo.updateProductStock(variantId, stock)
}

func (s *Service) FilterAllProducts(limit int, page int, searchString string, categoryFilter string, brandFilter string) ([]models.AllProducts, error) {
	return s.repo.filterAllProducts(limit, page, searchString, categoryFilter, brandFilter)
}

func (s *Service) FilterAllProductsCount(searchString string, categoryFilter string, brandFilter string) (int64, error) {
	return s.repo.filterAllProductsCount(searchString, categoryFilter, brandFilter)
}

func (s *Service) GenerateOrderIdByCart(price int, userId string, addressId string) (string, error) {
	return s.repo.generateOrderIdByCart(price, userId, addressId)
}

func (s *Service) AddProductsInOrder(orderedProducts []models.ProductOrdered) error {
	return s.repo.addProductsInOrder(orderedProducts)
}

func (s *Service) GetAllOffers() ([]models.Offer, error) {
	return s.repo.getAllOffers()
}

func (s *Service) AddAddress(userID string, newAddress *models.Address) (string, error) {
	return s.repo.addAddress(userID, newAddress)
}
