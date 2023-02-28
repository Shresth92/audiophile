package dbHelpers

import (
	"github.com/Shresth92/audiophile/database"
	"github.com/Shresth92/audiophile/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func GenerateOrderIdByCart(tx *gorm.DB, price int, userId uuid.UUID, addressId uuid.UUID) (uuid.UUID, error) {
	orderId := uuid.New()
	order := models.Orders{
		Id:             orderId,
		UserId:         userId,
		AddressId:      addressId,
		Cost:           price,
		DeliveryStatus: models.OnTheWay,
	}
	err := tx.
		Model(&models.Orders{}).
		Create(&order).
		Error
	return orderId, err
}

func AddProductsInOrder(tx *gorm.DB, OrderedProducts []models.ProductOrdered) error {
	err := tx.
		Model(&models.ProductOrdered{}).
		Create(&OrderedProducts).
		Error
	return err
}

func GetAllOffers() ([]models.Offer, error) {
	var offers []models.Offer
	err := database.Db.
		Model(&models.Offer{}).
		Where("archived_at is null").
		Scan(&offers).
		Error
	return offers, err
}
