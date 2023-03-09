package models

import (
	"time"
)

type (
	Category struct {
		Id           string    `json:"id" gorm:"column:id;index"`
		CategoryName string    `json:"categoryName" gorm:"column:category_name"`
		CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt    time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt   time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	Brand struct {
		Id         string    `json:"id" gorm:"column:id;primaryKey;index"`
		BrandName  string    `json:"brandName" gorm:"column:brand_name"`
		CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	Product struct {
		Id          string    `json:"id" gorm:"column:id;primaryKey;index"`
		ProductName string    `json:"productName" gorm:"column:product_name"`
		ModelName   string    `json:"modelName" gorm:"column:model_name"`
		BrandId     string    `json:"brandId" gorm:"column:brand_id"`
		Brand       Brand     `gorm:"foreignKey:BrandId"`
		CategoryId  string    `json:"categoryId"`
		Category    Category  `gorm:"column:category_id;foreignKey:CategoryId"`
		Return      int       `json:"return" gorm:"column:return"`
		Warranty    int       `json:"warranty"  gorm:"column:warranty"`
		Wireless    bool      `json:"wireless"  gorm:"column:wireless"`
		CreatedAt   time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt   time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt  time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	Variants struct {
		Id         string    `json:"id" gorm:"column:id;primaryKey;index"`
		ProductId  string    `json:"productId"`
		Product    Product   `gorm:"column:product_id;foreignKey:ProductId"`
		Colour     string    `json:"colour" gorm:"column:colour"`
		Price      int       `json:"price" gorm:"column:price"`
		Stock      int       `json:"stock" gorm:"column:stock"`
		CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	Offer struct {
		Id          string    `json:"id" gorm:"column:id;primaryKey;index"`
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
		Id         string    `json:"id" gorm:"column:id;primaryKey;index"`
		BucketName string    `json:"bucketName" gorm:"column:bucket_name"`
		Path       string    `json:"path" gorm:"column:path"`
		CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	VariantImages struct {
		Id         string    `json:"id" gorm:"column:id;primaryKey;index"`
		ImageId    string    `json:"imageId"`
		Images     Images    `gorm:"column:image_id;foreignKey:ImageId"`
		VariantId  string    `json:"variantId"`
		Variant    Variants  `gorm:"column:variant_id;foreignKey:VariantId"`
		CreatedAt  time.Time `json:"createdAt" gorm:"column:created_at;default:current_timestamp"`
		UpdatedAt  time.Time `json:"updatedAt" gorm:"column:updated_at;default:current_timestamp"`
		ArchivedAt time.Time `json:"archivedAt" gorm:"column:archived_at;default:null"`
	}

	AllProducts struct {
		VariantID    string `json:"variantId"`
		ProductName  string `json:"productName"`
		ModelName    string `json:"modelName"`
		BrandName    string `json:"brandName"`
		CategoryName string `json:"categoryName"`
		Return       int    `json:"return"`
		Warranty     int    `json:"warranty"`
		Wireless     bool   `json:"wireless"`
		Colour       string `json:"colour"`
		Price        int    `json:"price"`
		Stock        int    `json:"stock"`
		BucketName   string `json:"bucketName"`
		Path         string `json:"path"`
	}

	ProductBody struct {
		ProductName string   `json:"productName"`
		ModelName   string   `json:"modelName"`
		BrandID     string   `json:"brandId"`
		CategoryID  string   `json:"categoryId"`
		Return      int      `json:"return"`
		Warranty    int      `json:"warranty"`
		Wireless    bool     `json:"wireless"`
		Colour      string   `json:"colour"`
		Price       int      `json:"price"`
		Stock       int      `json:"stock"`
		ImageIds    []string `json:"imageIds"`
	}

	VariantBody struct {
		Colour   string   `json:"colour"`
		Price    int      `json:"price"`
		Stock    int      `json:"stock"`
		ImageIds []string `json:"imageIds"`
	}

	Products struct {
		Id           string   `json:"variantId"`
		ProductName  string   `json:"productName"`
		ModelName    string   `json:"modelName"`
		BrandName    string   `json:"brandName"`
		CategoryName string   `json:"categoryName"`
		Return       int      `json:"return"`
		Warranty     int      `json:"warranty"`
		Wireless     bool     `json:"wireless"`
		Colour       string   `json:"colour"`
		Price        int      `json:"price"`
		Stock        int      `json:"stock"`
		ImageLinks   []string `json:"imageLinks"`
	}
)
