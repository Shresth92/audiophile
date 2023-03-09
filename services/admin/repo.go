package admin

import (
	"github.com/Shresth92/audiophile/internal"
	"github.com/Shresth92/audiophile/models"
	"github.com/google/uuid"
	"time"
)

type repository struct {
	*internal.Database
}

func newCartRepository(db *internal.Database) *repository {
	return &repository{Database: db}
}

func (r *repository) uploadImageFirebase(bucket string, imagePath string) (string, error) {
	imageId := uuid.New().String()
	image := models.Images{
		Id:         imageId,
		BucketName: bucket,
		Path:       imagePath,
	}
	err := r.Database.DB.
		Model(&models.Images{}).
		Create(&image).
		Error
	return imageId, err
}

func (r *repository) createProduct(product *models.ProductBody) (string, error) {
	productId := uuid.New().String()
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
	err := r.Database.DB.
		Model(&models.Product{}).
		Create(&productInstance).
		Error
	return productId, err
}

func (r *repository) createCategory(categoryName string) error {
	categoryId := uuid.New().String()
	category := models.Category{
		Id:           categoryId,
		CategoryName: categoryName,
	}
	err := r.Database.DB.
		Model(&models.Category{}).
		Create(&category).
		Error
	return err
}

func (r *repository) createBrand(brandName string) error {
	brandId := uuid.New().String()
	brand := models.Brand{
		Id:        brandId,
		BrandName: brandName,
	}
	err := r.Database.DB.
		Model(&models.Brand{}).
		Create(&brand).
		Error
	return err
}

func (r *repository) createOffer(newOffer *models.Offer) error {
	offerId := uuid.New().String()
	offer := models.Offer{
		Id:          offerId,
		OfferName:   newOffer.OfferName,
		Percent:     newOffer.Percent,
		MaxDiscount: newOffer.MaxDiscount,
		CouponCode:  newOffer.CouponCode,
		Validity:    newOffer.Validity,
		Description: newOffer.Description,
	}
	err := r.Database.DB.
		Model(&models.Offer{}).
		Create(&offer).
		Error
	return err
}

func (r *repository) createVariant(productId string, colour string, stock int, price int) (string, error) {
	variantId := uuid.New().String()
	variant := models.Variants{
		Id:        variantId,
		ProductId: productId,
		Colour:    colour,
		Stock:     stock,
		Price:     price,
	}
	err := r.Database.DB.
		Model(&models.Variants{}).
		Create(&variant).Error
	return variantId, err
}

func (r *repository) uploadVariantImages(variantId string, imageIds []string) error {
	var variantImageArray []models.VariantImages
	for _, imageId := range imageIds {
		variantImageId := uuid.New().String()
		imageStruct := models.VariantImages{
			Id:        variantImageId,
			VariantId: variantId,
			ImageId:   imageId,
		}
		variantImageArray = append(variantImageArray, imageStruct)
	}

	err := r.Database.DB.
		Model(&models.VariantImages{}).
		Create(&variantImageArray).
		Error
	return err
}

func (r *repository) deleteVariant(productId string, variantId string) error {
	err := r.Database.DB.
		Model(&models.Variants{}).
		Where("id = ? and product_id = ? and archived_at is null", variantId, productId).
		Update("archived_at", time.Now()).
		Error
	return err
}

func (r *repository) deleteProduct(productId string) error {
	err := r.Database.DB.
		Model(&models.Product{}).
		Where("id = ? and archived_at is null", productId).
		Update("archived_at", time.Now()).
		Error
	return err
}

func (r *repository) deleteCategory(categoryId string) error {
	err := r.Database.DB.
		Model(&models.Category{}).
		Where("id = ? and archived_at is null", categoryId).
		Update("archived_at", time.Now()).
		Error
	return err
}

func (r *repository) deleteBrand(brandId string) error {
	err := r.Database.DB.
		Model(&models.Brand{}).
		Where("id = ? and archived_at is null", brandId).
		Update("archived_at", time.Now()).
		Error
	return err
}

func (r *repository) updateProduct(productId string, productDetails *models.Product) error {
	product := models.Product{
		ProductName: productDetails.ProductName,
		ModelName:   productDetails.ModelName,
		Return:      productDetails.Return,
		Warranty:    productDetails.Warranty,
		Wireless:    productDetails.Wireless,
		UpdatedAt:   time.Now(),
	}
	err := r.Database.DB.
		Model(&models.Product{}).
		Where("id =  ?", productId).
		Updates(&product).
		Error
	return err
}

func (r *repository) updateVariant(productId string, variantId string, variantDetails *models.Variants) error {
	variant := models.Variants{
		ProductId: productId,
		Colour:    variantDetails.Colour,
		Price:     variantDetails.Price,
		Stock:     variantDetails.Stock,
		UpdatedAt: time.Now(),
	}
	err := r.Database.DB.
		Model(&models.Variants{}).
		Where("id =  ?", variantId).
		Updates(&variant).
		Error
	return err
}

func (r *repository) updateCategory(categoryId string, categoryName string) error {
	category := models.Category{
		CategoryName: categoryName,
		UpdatedAt:    time.Now(),
	}
	err := r.Database.DB.
		Model(&models.Category{}).
		Where("id =  ?", categoryId).
		Updates(&category).
		Error
	return err
}

func (r *repository) updateBrand(brandId string, brandName string) error {
	brand := models.Brand{
		BrandName: brandName,
		UpdatedAt: time.Now(),
	}
	err := r.Database.DB.
		Model(&models.Brand{}).
		Where("id =  ?", brandId).
		Updates(&brand).
		Error
	return err
}

func (r *repository) getAllUsers(limit int, page int) ([]models.Users, error) {
	var users []models.Users
	err := r.Database.DB.
		Model(&models.Users{}).
		Preload("Address").
		Where("archived_at is null").
		Limit(limit).
		Offset(limit * (page - 1)).
		Find(&users).
		Error
	return users, err
}

func (r *repository) usersCount() (int64, error) {
	var count int64
	err := r.Database.DB.
		Model(&models.Users{}).
		Where("archived_at is null").
		Count(&count).
		Error
	return count, err
}

func (r *repository) getAllBrands(limit int, page int) ([]models.Brand, error) {
	var brands []models.Brand
	err := r.Database.DB.
		Model(&models.Brand{}).
		Where("archived_at is null").
		Limit(limit).
		Offset(limit * (page - 1)).
		Scan(&brands).
		Error
	return brands, err
}

func (r *repository) getBrandsCount() (int64, error) {
	var count int64
	err := r.Database.DB.
		Model(&models.Brand{}).
		Where("archived_at is null").
		Count(&count).
		Error
	return count, err
}

func (r *repository) getAllCategory(limit int, page int) ([]models.Category, error) {
	var categories []models.Category
	err := r.Database.DB.
		Model(&models.Category{}).
		Where("archived_at is null").
		Limit(limit).
		Offset(limit * (page - 1)).
		Scan(&categories).
		Error
	return categories, err
}

func (r *repository) getCategoryCount() (int64, error) {
	var count int64
	err := r.Database.DB.
		Model(&models.Category{}).
		Where("archived_at is null").
		Count(&count).
		Error
	return count, err
}

func (r *repository) changeUserRole(userId string, adminId string) error {
	roleId := uuid.New().String()
	role := models.UserRole{
		Id:        roleId,
		UserId:    userId,
		Role:      models.Admin,
		CreatedBy: adminId,
	}
	err := r.Database.DB.
		Model(&models.UserRole{}).
		Create(&role).
		Error
	return err
}
