package models

import (
	"time"
)

type (
	UserCart struct {
		Id         string    `json:"id" gorm:"column:id;primaryKey;index"`
		VariantId  string    `json:"variantId"`
		Variant    Variants  `gorm:"foreignKey:VariantId"`
		UserId     string    `json:"userId"`
		User       Users     `gorm:"foreignKey:UserId"`
		Count      int       `json:"count" gorm:"column:count"`
		CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}
)
