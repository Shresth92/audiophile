package dbHelpers

import (
	"github.com/Shresth92/audiophile/database"
	"github.com/Shresth92/audiophile/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func AddProductToCart(userId uuid.UUID, variantId uuid.UUID, count int) (uuid.UUID, error) {
	cartId := uuid.New()
	cart := models.UserCart{
		Id:        cartId,
		VariantId: variantId,
		UserId:    userId,
		Count:     count,
	}
	err := database.Db.
		Model(&models.UserCart{}).
		Create(&cart).
		Error
	return cartId, err
}

func GetCartProducts(userId uuid.UUID) ([]models.UserCart, error) {
	var cartItems []models.UserCart
	err := database.Db.
		Model(&models.UserCart{}).
		Where("user_id=? and archived_at is null", userId).
		Scan(&cartItems).
		Error
	return cartItems, err
}

func UpdateProductCountInCart(userId uuid.UUID, variantId uuid.UUID, count bool) error {
	var expression string
	if count {
		expression = "count + ?"
	} else {
		expression = "count - ?"
	}
	err := database.Db.
		Model(&models.UserCart{}).
		Where("user_id = ? and variant_id = ?", userId, variantId).
		UpdateColumn("count", gorm.Expr(expression, 1)).
		Error
	return err
}

func RemoveCartProduct(userId uuid.UUID, variantId uuid.UUID) error {
	err := database.Db.
		Model(&models.UserCart{}).
		Where("user_id = ? and variant_id = ?", userId, variantId).
		Update("archived_at", time.Now()).
		Error
	return err
}

func DeleteCart(tx *gorm.DB, userId uuid.UUID) error {
	err := tx.
		Model(&models.UserCart{}).
		Where("user_id = ?", userId).
		Update("archived_at", time.Now()).
		Error
	return err
}
