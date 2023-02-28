package models

import (
	"github.com/google/uuid"
	"time"
)

type Location struct {
	X, Y float64
}

type Roles string

const (
	Admin Roles = "admin"
	User  Roles = "user"
)

type (
	Users struct {
		Id         uuid.UUID `json:"id" gorm:"column:id;primaryKey;index"`
		Email      string    `json:"email" gorm:"column:email;size:255;index:unique_email,unique,where:archived_at is not null"`
		Password   string    `json:"password" gorm:"column:password;size:255"`
		Address    []Address `gorm:"foreignKey:UserId;references:Id"`
		CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archived_at" gorm:"column:archived_at;default:null"`
	}

	UserRole struct {
		Id         uuid.UUID `json:"id" gorm:"column:id;primaryKey;index"`
		UserId     uuid.UUID `json:"user_id"`
		User       Users     `gorm:"foreignKey:UserId"`
		Role       Roles     `gorm:"type:role_type"`
		CreatedBy  uuid.UUID `json:"created_by" gorm:"created_by"`
		Created    Users     `gorm:"foreignKey:UserId"`
		CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archived_at" gorm:"column:archived_at;default:null"`
	}

	Address struct {
		Id         uuid.UUID `json:"id" gorm:"column:id;primaryKey;index"`
		UserId     uuid.UUID `json:"user_id"`
		Area       string    `json:"area" gorm:"column:area"`
		City       string    `json:"city" gorm:"column:city"`
		State      string    `json:"state" gorm:"column:state"`
		ZipCode    string    `json:"zip_code" gorm:"column:zipcode;type:varchar(6)"`
		Contact    string    `json:"contact" gorm:"column:contact"`
		LatLong    string    `json:"lat_long" gorm:"column:lat_long;type:Point"`
		CreatedAt  time.Time `json:"created_at" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updated_at" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archived_at" gorm:"column:archived_at;default:null"`
	}

	Session struct {
		ID        uuid.UUID `json:"id" gorm:"column:id;primaryKey;index"`
		UserId    uuid.UUID `json:"user_id"`
		StartedAt time.Time `json:"started_at" gorm:"column:started_at;default:current_timestamp"`
		EndedAt   time.Time `json:"ended_at" gorm:"default:null"`
		User      Users     `gorm:"foreignKey:UserId"`
	}
)
