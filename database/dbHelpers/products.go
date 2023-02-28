package dbHelpers

import (
	"errors"
	"github.com/Shresth92/audiophile/database"
	"github.com/Shresth92/audiophile/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func GetProduct(productId uuid.UUID) ([]models.AllProducts, error) {
	var product []models.AllProducts
	err := database.Db.Table("products").
		Joins("join brands b on products.brand_id=b.id").
		Joins("join categories c on products.category_id=c.id").
		Joins("join variants v on v.product_id=products.id").
		Joins("left join variant_images vi on v.id=vi.variant_id").
		Joins("join images i on i.id=vi.image_id").
		Where("products.id = ?", productId).Find(&product).Error
	return product, err
}

func GetTotalProductCost(variantIds []uuid.UUID) (int, error) {
	var costs []int
	totalCost := 0
	err := database.Db.Model(&models.Variants{}).Where("id IN ?", variantIds).Pluck("price", &costs).Error
	for _, cost := range costs {
		totalCost += cost
	}
	return totalCost, err
}

func PriceAfterDiscount(price int, couponCode string) (int, error) {
	discountedPrice := price
	discount := models.Offer{}
	err := database.Db.Model(&models.Offer{}).Where("coupon_code = ?", couponCode).Find(&discount).Error
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

func FilterMyOrders(userId uuid.UUID, ProductStatus models.DeliveryStatus, limit int, page int) ([]models.Orders, error) {
	var orders []models.Orders
	err := database.Db.
		Preload("ProductOrdered").
		Where("user_id = ? and delivery_status = ?", userId, ProductStatus).
		Limit(limit).
		Offset(limit * (page - 1)).
		Find(&orders).Error
	return orders, err
}

func CountFilterMyOrders(userId uuid.UUID, ProductStatus models.DeliveryStatus) (int64, error) {
	var count int64
	err := database.Db.Where("user_id = ? and delivery_status = ?", userId, ProductStatus).Count(&count).Error
	return count, err
}

func UpdateProductStock(tx *gorm.DB, variantId uuid.UUID, stock int) error {
	err := tx.
		Model(&models.Variants{}).
		Where("id = ? and stock >= ? and archived_at is null", variantId, stock).
		UpdateColumn("stock", gorm.Expr("stock - ?", stock)).
		Error
	return err
}

func FilterAllProducts(limit int, page int, searchString string, categoryFilter string, brandFilter string) ([]models.AllProducts, error) {
	var product []models.AllProducts
	if searchString != "" {
		database.Db = database.Db.Where("products.product_name ilike ?", searchString)
	}
	if categoryFilter != "" {
		database.Db = database.Db.Or("c.category_name = ?", categoryFilter)
	}
	if brandFilter != "" {
		database.Db = database.Db.Or("b.brand_name = ?", brandFilter)
	}
	dbQuery := database.Db.Table("products").
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

func FilterAllProductsCount(searchString string, categoryFilter string, brandFilter string) (int64, error) {
	var count int64
	if searchString != "" {
		database.Db = database.Db.Where("products.product_name ilike ?", searchString)
	}
	if categoryFilter != "" {
		database.Db = database.Db.Or("c.category_name = ?", categoryFilter)
	}
	if brandFilter != "" {
		database.Db = database.Db.Or("b.brand_name = ?", brandFilter)
	}
	err := database.Db.Table("products").
		Joins("join brands b on products.brand_id=b.id").
		Joins("join categories c on products.category_id=c.id").
		Joins("join variants v on v.product_id=products.id").
		Joins("left join variant_images vi on v.id=vi.variant_id").
		Joins("join images i on i.id=vi.image_id").
		Count(&count).
		Error
	return count, err
}
