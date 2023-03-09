package models

import (
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
		ID             string           `json:"id" gorm:"column:id;primaryKey;index"`
		UserID         string           `json:"userId"`
		User           Users            `gorm:"foreignKey:UserID"`
		OrderedAt      time.Time        `json:"orderedAt" gorm:"column:ordered_at;default:current_timestamp"`
		DeliveredAt    time.Time        `json:"deliveredAt" gorm:"column:delivered_at"`
		AddressId      string           `json:"addressId"`
		ProductOrdered []ProductOrdered `gorm:"foreignKey:OrderId;references:ID"`
		Address        Address          `gorm:"foreignKey:AddressId"`
		Cost           int              `json:"cost" gorm:"column:cost"`
		DeliveryStatus DeliveryStatus   `json:"deliveryStatus" gorm:"column:delivery_status;type:delivery_status"`
		CreatedAt      time.Time        `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt      time.Time        `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt     time.Time        `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	ProductOrdered struct {
		ID         string    `json:"id" gorm:"column:id;primaryKey;index"`
		VariantId  string    `json:"variant_id"`
		Variant    Variants  `gorm:"foreignKey:VariantId"`
		Quantity   int       `json:"quantity" gorm:"column:quantity"`
		OrderId    string    `json:"orderId"`
		CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}
)
