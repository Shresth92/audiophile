package models

import (
	"github.com/google/uuid"
	"time"
)

type DeliveryStatus string

const (
	OnTheWay  DeliveryStatus = "onTheWay"
	Delivered DeliveryStatus = "delivered"
	Canceled  DeliveryStatus = "canceled"
	Return    DeliveryStatus = "return"
)

type (
	Orders struct {
		Id             uuid.UUID        `json:"id" gorm:"column:id;primaryKey;index"`
		UserId         uuid.UUID        `json:"userId"`
		User           Users            `gorm:"foreignKey:UserId"`
		OrderedAt      time.Time        `json:"orderedAt" gorm:"column:ordered_at;default:current_timestamp"`
		DeliveredAt    time.Time        `json:"deliveredAt" gorm:"column:delivered_at"`
		AddressId      uuid.UUID        `json:"addressId"`
		ProductOrdered []ProductOrdered `gorm:"foreignKey:OrderId;references:Id"`
		Address        Address          `gorm:"column:addressId;foreignKey:AddressId"`
		Cost           int              `json:"cost" gorm:"column:cost"`
		DeliveryStatus DeliveryStatus   `json:"deliveryStatus" gorm:"column:delivery_status;type:delivery_status"`
		CreatedAt      time.Time        `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt      time.Time        `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt     time.Time        `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	ProductOrdered struct {
		Id         uuid.UUID `json:"id" gorm:"column:id;primaryKey;index"`
		VariantId  uuid.UUID `json:"variant_id"`
		Variant    Variants  `gorm:"foreignKey:VariantId"`
		Quantity   int       `json:"quantity" gorm:"column:quantity"`
		OrderId    uuid.UUID `json:"orderId"`
		CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}
)
