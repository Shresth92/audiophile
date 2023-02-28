package dbHelpers

import (
	"github.com/Shresth92/audiophile/database"
	"github.com/Shresth92/audiophile/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

func UploadImageFirebase(bucket string, imagePath string) (uuid.UUID, error) {
	imageId := uuid.New()
	image := models.Images{
		Id:         imageId,
		BucketName: bucket,
		Path:       imagePath,
	}
	err := database.Db.
		Model(&models.Images{}).
		Create(&image).
		Error
	return imageId, err
}

func CreateProduct(tx *gorm.DB, product *models.ProductBody) (uuid.UUID, error) {
	productId := uuid.New()
	productInstance := models.Product{
		Id:          productId,
		ProductName: product.ProductName,
		ModelName:   product.ModelName,
		BrandId:     product.BrandID,
		CategoryId:  product.CategoryID,
		Return:      product.Return,
		Warranty:    product.Warranty,
		Wireless:    product.Wireless,
	}
	err := tx.
		Model(&models.Product{}).
		Create(&productInstance).
		Error
	return productId, err
}

func CreateCategory(categoryName string) (uuid.UUID, error) {
	categoryId := uuid.New()
	category := models.Category{
		Id:           categoryId,
		CategoryName: categoryName,
	}
	err := database.Db.
		Model(&models.Category{}).
		Create(&category).
		Error
	return categoryId, err
}

func CreateBrand(brandName string) (uuid.UUID, error) {
	brandId := uuid.New()
	brand := models.Brand{
		Id:        brandId,
		BrandName: brandName,
	}
	err := database.Db.
		Model(&models.Brand{}).
		Create(&brand).
		Error
	return brandId, err
}

func CreateOffer(newOffer *models.Offer) (uuid.UUID, error) {
	offerId := uuid.New()
	offer := models.Offer{
		Id:          offerId,
		OfferName:   newOffer.OfferName,
		Percent:     newOffer.Percent,
		MaxDiscount: newOffer.MaxDiscount,
		CouponCode:  newOffer.CouponCode,
		Validity:    newOffer.Validity,
		Description: newOffer.Description,
	}
	err := database.Db.
		Model(&models.Offer{}).
		Create(&offer).
		Error
	return offerId, err
}

func CreateVariant(tx *gorm.DB, productId uuid.UUID, colour string, stock int, price int) (uuid.UUID, error) {
	variantId := uuid.New()
	variant := models.Variants{
		Id:        variantId,
		ProductId: productId,
		Colour:    colour,
		Stock:     stock,
		Price:     price,
	}
	err := tx.
		Model(&models.Variants{}).
		Create(&variant).Error
	return variantId, err
}

func UploadVariantImages(tx *gorm.DB, variantId uuid.UUID, imageIds []string) ([]uuid.UUID, error) {
	var variantImageArray []models.VariantImages
	var variantImageIds []uuid.UUID
	for _, imageId := range imageIds {
		variantImageId := uuid.New()
		variantImageIds = append(variantImageIds, variantImageId)
		imageStruct := models.VariantImages{
			Id:        variantImageId,
			VariantId: variantId,
			ImageId:   uuid.MustParse(imageId),
		}
		variantImageArray = append(variantImageArray, imageStruct)
	}

	err := tx.
		Model(&models.VariantImages{}).
		Create(&variantImageArray).
		Error
	return variantImageIds, err
}

func DeleteVariant(productId uuid.UUID, variantId uuid.UUID) error {
	err := database.Db.
		Model(&models.Variants{}).
		Where("id = ? and product_id = ? and archived_at is null", variantId, productId).
		Update("archived_at", time.Now()).
		Error
	return err
}

func DeleteProduct(productId uuid.UUID) error {
	err := database.Db.
		Model(&models.Product{}).
		Where("id = ? and archived_at is null", productId).
		Update("archived_at", time.Now()).
		Error
	return err
}

func DeleteCategory(categoryId uuid.UUID) error {
	err := database.Db.
		Model(&models.Category{}).
		Where("id = ? and archived_at is null", categoryId).
		Update("archived_at", time.Now()).
		Error
	return err
}

func DeleteBrand(brandId uuid.UUID) error {
	err := database.Db.
		Model(&models.Brand{}).
		Where("id = ? and archived_at is null", brandId).
		Update("archived_at", time.Now()).
		Error
	return err
}

func UpdateProduct(productId uuid.UUID, productName string, modelName string, returnDays int, warranty int, wireless bool) error {
	product := models.Product{
		ProductName: productName,
		ModelName:   modelName,
		Return:      returnDays,
		Warranty:    warranty,
		Wireless:    wireless,
		UpdatedAt:   time.Now(),
	}
	err := database.Db.
		Model(&models.Product{}).
		Where("id =  ?", productId).
		Updates(&product).
		Error
	return err
}

func UpdateVariant(productId uuid.UUID, variantId uuid.UUID, colour string, price int, stock int) error {
	variant := models.Variants{
		ProductId: productId,
		Colour:    colour,
		Price:     price,
		Stock:     stock,
		UpdatedAt: time.Now(),
	}
	err := database.Db.
		Model(&models.Variants{}).
		Where("id =  ?", variantId).
		Updates(&variant).
		Error
	return err
}

func UpdateCategory(categoryId uuid.UUID, categoryName string) error {
	category := models.Category{
		CategoryName: categoryName,
		UpdatedAt:    time.Now(),
	}
	err := database.Db.
		Model(&models.Category{}).
		Where("id =  ?", categoryId).
		Updates(&category).
		Error
	return err
}

func UpdateBrand(brandId uuid.UUID, brandName string) error {
	brand := models.Brand{
		BrandName: brandName,
		UpdatedAt: time.Now(),
	}
	err := database.Db.
		Model(&models.Brand{}).
		Where("id =  ?", brandId).
		Updates(&brand).
		Error
	return err
}

func GetAllUsers(limit int, page int) ([]models.Users, error) {
	var users []models.Users
	err := database.Db.
		Model(&models.Users{}).
		Preload("Address").
		Where("archived_at is null").
		Limit(limit).
		Offset(limit * (page - 1)).
		Find(&users).
		Error
	return users, err
}

func UsersCount() (int64, error) {
	var count int64
	err := database.Db.
		Model(&models.Users{}).
		Where("archived_at is null").
		Count(&count).
		Error
	return count, err
}

func GetAllBrands(limit int, page int) ([]models.Brand, error) {
	var brands []models.Brand
	err := database.Db.
		Model(&models.Brand{}).
		Where("archived_at is null").
		Limit(limit).
		Offset(limit * (page - 1)).
		Scan(&brands).
		Error
	return brands, err
}

func GetBrandsCount() (int64, error) {
	var count int64
	err := database.Db.
		Model(&models.Brand{}).
		Where("archived_at is null").
		Count(&count).
		Error
	return count, err
}

func GetAllCategory(limit int, page int) ([]models.Category, error) {
	var categories []models.Category
	err := database.Db.
		Model(&models.Category{}).
		Where("archived_at is null").
		Limit(limit).
		Offset(limit * (page - 1)).
		Scan(&categories).
		Error
	return categories, err
}

func GetCategoryCount() (int64, error) {
	var count int64
	err := database.Db.
		Model(&models.Category{}).
		Where("archived_at is null").
		Count(&count).
		Error
	return count, err
}

func ChangeUserRole(userId uuid.UUID, adminId uuid.UUID) error {
	roleId := uuid.New()
	role := models.UserRole{
		Id:        roleId,
		UserId:    userId,
		Role:      models.Admin,
		CreatedBy: adminId,
	}
	err := database.Db.
		Model(&models.UserRole{}).
		Create(&role).
		Error
	return err
}
