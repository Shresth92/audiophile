package user

import (
	"errors"
	"github.com/Shresth92/audiophile/internal"
	"github.com/Shresth92/audiophile/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type repository struct {
	*internal.Database
}

func newUserRepository(db *internal.Database) *repository {
	return &repository{Database: db}
}

func (r *repository) addProductToCart(userId string, variantId string, count int) error {
	cartId := uuid.New().String()
	cart := models.UserCart{
		Id:        cartId,
		VariantId: variantId,
		UserId:    userId,
		Count:     count,
	}
	err := r.Database.DB.
		Model(&models.UserCart{}).
		Create(&cart).
		Error
	return err
}

func (r *repository) getCartProducts(userId string) ([]models.UserCart, error) {
	var cartItems []models.UserCart
	err := r.Database.DB.
		Model(&models.UserCart{}).
		Where("user_id=? and archived_at is null", userId).
		Scan(&cartItems).
		Error
	return cartItems, err
}

func (r *repository) updateProductCountInCart(userId string, variantId string, count bool) error {
	var expression string
	if count {
		expression = "count + ?"
	} else {
		expression = "count - ?"
	}
	err := r.Database.DB.
		Model(&models.UserCart{}).
		Where("user_id = ? and variant_id = ?", userId, variantId).
		UpdateColumn("count", gorm.Expr(expression, 1)).
		Error
	return err
}

func (r *repository) removeCartProduct(userId string, variantId string) error {
	err := r.Database.DB.
		Model(&models.UserCart{}).
		Where("user_id = ? and variant_id = ?", userId, variantId).
		Update("archived_at", time.Now()).
		Error
	return err
}

func (r *repository) deleteCart(userId string) error {
	err := r.Database.DB.
		Model(&models.UserCart{}).
		Where("user_id = ?", userId).
		Update("archived_at", time.Now()).
		Error
	return err
}

func (r *repository) getProduct(productId string) ([]models.AllProducts, error) {
	var product []models.AllProducts
	err := r.Database.DB.Table("products").
		Joins("join brands b on products.brand_id=b.id").
		Joins("join categories c on products.category_id=c.id").
		Joins("join variants v on v.product_id=products.id").
		Joins("left join variant_images vi on v.id=vi.variant_id").
		Joins("join images i on i.id=vi.image_id").
		Where("products.id = ?", productId).Find(&product).Error
	return product, err
}

func (r *repository) getTotalProductCost(variantIds []string) (int, error) {
	var costs []int
	totalCost := 0
	err := r.Database.DB.Model(&models.Variants{}).Where("id IN ?", variantIds).Pluck("price", &costs).Error
	for _, cost := range costs {
		totalCost += cost
	}
	return totalCost, err
}

func (r *repository) priceAfterDiscount(price int, couponCode string) (int, error) {
	discountedPrice := price
	discount := models.Offer{}
	err := r.Database.DB.Model(&models.Offer{}).Where("coupon_code = ?", couponCode).Find(&discount).Error
	if discount.Validity.Before(time.Now()) {
		return price, errors.New("coupon is expired")
	}
	if discount.Percent != 0 {
		var prodDiscount int
		prodDiscount = discount.Percent * 100 / price
		if prodDiscount > discount.MaxDiscount {
			prodDiscount = discount.MaxDiscount
		}
		discountedPrice = price - prodDiscount
	}
	return discountedPrice, err
}

func (r *repository) filterMyOrders(userId string, ProductStatus models.DeliveryStatus, limit int, page int) ([]models.Orders, error) {
	var orders []models.Orders
	err := r.Database.DB.
		Preload("ProductOrdered").
		Where("user_id = ? and delivery_status = ?", userId, ProductStatus).
		Limit(limit).
		Offset(limit * (page - 1)).
		Find(&orders).Error
	return orders, err
}

func (r *repository) countFilterMyOrders(userId string, ProductStatus models.DeliveryStatus) (int64, error) {
	var count int64
	err := r.Database.DB.Where("user_id = ? and delivery_status = ?", userId, ProductStatus).Count(&count).Error
	return count, err
}

func (r *repository) updateProductStock(variantId string, stock int) error {
	err := r.Database.DB.
		Model(&models.Variants{}).
		Where("id = ? and stock >= ? and archived_at is null", variantId, stock).
		UpdateColumn("stock", gorm.Expr("stock - ?", stock)).
		Error
	return err
}

func (r *repository) filterAllProducts(limit int, page int, searchString string, categoryFilter string, brandFilter string) ([]models.AllProducts, error) {
	var product []models.AllProducts
	if searchString != "" {
		r.Database.DB = r.Database.DB.Where("products.product_name ilike ?", searchString)
	}
	if categoryFilter != "" {
		r.Database.DB = r.Database.DB.Or("c.category_name = ?", categoryFilter)
	}
	if brandFilter != "" {
		r.Database.DB = r.Database.DB.Or("b.brand_name = ?", brandFilter)
	}
	dbQuery := r.Database.DB.Table("products").
		Joins("join brands b on products.brand_id=b.id").
		Joins("join categories c on products.category_id=c.id").
		Joins("join variants v on v.product_id=products.id").
		Joins("left join variant_images vi on v.id=vi.variant_id").
		Joins("join images i on i.id=vi.image_id").
		Limit(limit).
		Offset(limit * (page - 1)).
		Scan(&product)
	return product, dbQuery.Error
}

func (r *repository) filterAllProductsCount(searchString string, categoryFilter string, brandFilter string) (int64, error) {
	var count int64
	if searchString != "" {
		r.Database.DB = r.Database.DB.Where("products.product_name ilike ?", searchString)
	}
	if categoryFilter != "" {
		r.Database.DB = r.Database.DB.Or("c.category_name = ?", categoryFilter)
	}
	if brandFilter != "" {
		r.Database.DB = r.Database.DB.Or("b.brand_name = ?", brandFilter)
	}
	err := r.Database.DB.Table("products").
		Joins("join brands b on products.brand_id=b.id").
		Joins("join categories c on products.category_id=c.id").
		Joins("join variants v on v.product_id=products.id").
		Joins("left join variant_images vi on v.id=vi.variant_id").
		Joins("join images i on i.id=vi.image_id").
		Count(&count).
		Error
	return count, err
}

func (r *repository) generateOrderIdByCart(price int, userId string, addressId string) (string, error) {
	orderId := uuid.New().String()
	order := models.Orders{
		ID:             orderId,
		UserID:         userId,
		AddressId:      addressId,
		Cost:           price,
		DeliveryStatus: models.OnTheWay,
	}
	err := r.Database.DB.
		Model(&models.Orders{}).
		Create(&order).
		Error
	return orderId, err
}

func (r *repository) addProductsInOrder(OrderedProducts []models.ProductOrdered) error {
	err := r.Database.DB.
		Model(&models.ProductOrdered{}).
		Create(&OrderedProducts).
		Error
	return err
}

func (r *repository) getAllOffers() ([]models.Offer, error) {
	var offers []models.Offer
	err := r.Database.DB.
		Model(&models.Offer{}).
		Where("archived_at is null").
		Scan(&offers).
		Error
	return offers, err
}

func (r *repository) addAddress(userId string, newAddress *models.Address) (string, error) {
	addressId := uuid.New().String()
	address := models.Address{
		Id:      addressId,
		UserId:  userId,
		Area:    newAddress.Area,
		City:    newAddress.City,
		State:   newAddress.State,
		ZipCode: newAddress.ZipCode,
		Contact: newAddress.Contact,
		LatLong: newAddress.LatLong,
	}
	err := r.Database.DB.
		Model(&models.Address{}).
		Create(&address).
		Error
	return addressId, err
}
