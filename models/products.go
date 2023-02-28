package models

import (
	"github.com/google/uuid"
	"time"
)

type (
	Category struct {
		Id           uuid.UUID `json:"id" gorm:"column:id;index"`
		CategoryName string    `json:"categoryName" gorm:"column:category_name"`
		CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt   time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	Brand struct {
		Id         uuid.UUID `json:"id" gorm:"column:id;primaryKey;index"`
		BrandName  string    `json:"brandName" gorm:"column:brand_name"`
		CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	Product struct {
		Id          uuid.UUID `json:"id" gorm:"column:id;primaryKey;index"`
		ProductName string    `json:"productName" gorm:"column:product_name"`
		ModelName   string    `json:"modelName" gorm:"column:model_name"`
		BrandId     uuid.UUID `json:"brandId" gorm:"column:brand_id"`
		Brand       Brand     `gorm:"foreignKey:BrandId"`
		CategoryId  uuid.UUID `json:"categoryId"`
		Category    Category  `gorm:"column:category_id;foreignKey:CategoryId"`
		Return      int       `json:"return" gorm:"column:return"`
		Warranty    int       `json:"warranty"  gorm:"column:warranty"`
		Wireless    bool      `json:"wireless"  gorm:"column:wireless"`
		CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt  time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	Variants struct {
		Id         uuid.UUID `json:"id" gorm:"column:id;primaryKey;index"`
		ProductId  uuid.UUID `json:"productId"`
		Product    Product   `gorm:"column:product_id;foreignKey:ProductId"`
		Colour     string    `json:"colour" gorm:"column:colour"`
		Price      int       `json:"price" gorm:"column:price"`
		Stock      int       `json:"stock" gorm:"column:stock"`
		CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	Offer struct {
		Id          uuid.UUID `json:"id" gorm:"column:id;primaryKey;index"`
		OfferName   string    `json:"offerName" gorm:"column:offer_name"`
		Percent     int       `json:"percent" gorm:"column:percent"`
		MaxDiscount int       `json:"maxDiscount" gorm:"column:max_discount"`
		CouponCode  string    `json:"couponCode" gorm:"column:coupon_code"`
		Validity    time.Time `json:"validity"  gorm:"column:validity"`
		Description string    `json:"description" gorm:"column:description"`
		CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt  time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	Images struct {
		Id         uuid.UUID `json:"id" gorm:"column:id;primaryKey;index"`
		BucketName string    `json:"bucketName" gorm:"column:bucket_name"`
		Path       string    `json:"path" gorm:"column:path"`
		CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	VariantImages struct {
		Id         uuid.UUID `json:"id" gorm:"column:id;primaryKey;index"`
		ImageId    uuid.UUID `json:"imageId"`
		Images     Images    `gorm:"column:image_id;foreignKey:ImageId"`
		VariantId  uuid.UUID `json:"variantId"`
		Variant    Variants  `gorm:"column:variant_id;foreignKey:VariantId"`
		CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	AllProducts struct {
		Id           uuid.UUID `json:"variantId"`
		ProductName  string    `json:"productName"`
		ModelName    string    `json:"modelName"`
		BrandName    string    `json:"brandName"`
		CategoryName string    `json:"categoryName"`
		Return       int       `json:"return"`
		Warranty     int       `json:"warranty"`
		Wireless     bool      `json:"wireless"`
		Colour       string    `json:"colour"`
		Price        int       `json:"price"`
		Stock        int       `json:"stock"`
		BucketName   string    `json:"bucketName"`
		Path         string    `json:"path"`
	}

	ProductBody struct {
		ProductName string    `json:"productName"`
		ModelName   string    `json:"modelName"`
		BrandID     uuid.UUID `json:"brandId"`
		CategoryID  uuid.UUID `json:"categoryId"`
		Return      int       `json:"return"`
		Warranty    int       `json:"warranty"`
		Wireless    bool      `json:"wireless"`
		Colour      string    `json:"colour"`
		Price       int       `json:"price"`
		Stock       int       `json:"stock"`
		ImageIds    []string  `json:"imageIds"`
	}

	VariantBody struct {
		Colour   string   `json:"colour"`
		Price    int      `json:"price"`
		Stock    int      `json:"stock"`
		ImageIds []string `json:"imageIds"`
	}

	Products struct {
		Id           uuid.UUID `json:"variantId"`
		ProductName  string    `json:"productName"`
		ModelName    string    `json:"modelName"`
		BrandName    string    `json:"brandName"`
		CategoryName string    `json:"categoryName"`
		Return       int       `json:"return"`
		Warranty     int       `json:"warranty"`
		Wireless     bool      `json:"wireless"`
		Colour       string    `json:"colour"`
		Price        int       `json:"price"`
		Stock        int       `json:"stock"`
		ImageLinks   []string  `json:"imageLinks"`
	}
)
