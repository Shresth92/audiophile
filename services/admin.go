package services

import (
	"github.com/Shresth92/audiophile/models"
)

type AdminServices interface {
	UploadImageFirebase(bucket string, imagePath string) (string, error)
	CreateProduct(product *models.ProductBody) (string, error)
	CreateCategory(categoryName string) error
	CreateBrand(brandName string) error
	CreateOffer(newOffer *models.Offer) error
	CreateVariant(productId string, colour string, stock int, price int) (string, error)
	UploadVariantImages(variantId string, imageIds []string) error
	DeleteVariant(productId string, variantId string) error
	DeleteProduct(productId string) error
	DeleteCategory(categoryId string) error
	DeleteBrand(brandId string) error
	UpdateProduct(productId string, productDetails *models.Product) error
	UpdateVariant(productId string, variantId string, variantDetails *models.Variants) error
	UpdateCategory(categoryId string, categoryName string) error
	UpdateBrand(brandId string, brandName string) error
	GetAllUsers(limit int, page int) ([]models.Users, error)
	UsersCount() (int64, error)
	GetAllBrands(limit int, page int) ([]models.Brand, error)
	GetBrandsCount() (int64, error)
	GetAllCategory(limit int, page int) ([]models.Category, error)
	GetCategoryCount() (int64, error)
	ChangeUserRole(userId string, adminId string) error
}
