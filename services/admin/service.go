package admin

import (
	"github.com/Shresth92/audiophile/internal"
	"github.com/Shresth92/audiophile/models"
)

type Service struct {
	repo *repository
}

func NewAdminService(db *internal.Database) *Service {
	return &Service{repo: newCartRepository(db)}
}

func (s *Service) UploadImageFirebase(bucket string, imagePath string) (string, error) {
	return s.repo.uploadImageFirebase(bucket, imagePath)
}

func (s *Service) CreateProduct(product *models.ProductBody) (string, error) {
	return s.repo.createProduct(product)
}

func (s *Service) CreateCategory(categoryName string) error {
	return s.repo.createCategory(categoryName)
}

func (s *Service) CreateBrand(brandName string) error {
	return s.repo.createBrand(brandName)
}

func (s *Service) CreateOffer(newOffer *models.Offer) error {
	return s.repo.createOffer(newOffer)
}

func (s *Service) CreateVariant(productId string, colour string, stock int, price int) (string, error) {
	return s.repo.createVariant(productId, colour, stock, price)
}

func (s *Service) UploadVariantImages(variantId string, imageIds []string) error {
	return s.repo.uploadVariantImages(variantId, imageIds)
}

func (s *Service) DeleteVariant(productId string, variantId string) error {
	return s.repo.deleteVariant(productId, variantId)
}

func (s *Service) DeleteProduct(productId string) error {
	return s.repo.deleteProduct(productId)
}

func (s *Service) DeleteCategory(categoryId string) error {
	return s.repo.deleteCategory(categoryId)
}

func (s *Service) DeleteBrand(brandId string) error {
	return s.repo.deleteBrand(brandId)
}

func (s *Service) UpdateProduct(productId string, productDetails *models.Product) error {
	return s.repo.updateProduct(productId, productDetails)
}

func (s *Service) UpdateVariant(productId string, variantId string, variantDetails *models.Variants) error {
	return s.repo.updateVariant(productId, variantId, variantDetails)
}

func (s *Service) UpdateCategory(categoryId string, categoryName string) error {
	return s.repo.updateCategory(categoryId, categoryName)
}

func (s *Service) UpdateBrand(brandId string, brandName string) error {
	return s.repo.updateBrand(brandId, brandName)
}

func (s *Service) GetAllUsers(limit int, page int) ([]models.Users, error) {
	return s.repo.getAllUsers(limit, page)
}

func (s *Service) UsersCount() (int64, error) {
	return s.repo.usersCount()
}

func (s *Service) GetAllBrands(limit int, page int) ([]models.Brand, error) {
	return s.repo.getAllBrands(limit, page)
}

func (s *Service) GetBrandsCount() (int64, error) {
	return s.repo.getBrandsCount()
}

func (s *Service) GetAllCategory(limit int, page int) ([]models.Category, error) {
	return s.repo.getAllCategory(limit, page)
}

func (s *Service) GetCategoryCount() (int64, error) {
	return s.repo.getCategoryCount()
}

func (s *Service) ChangeUserRole(userId string, adminId string) error {
	return s.repo.changeUserRole(userId, adminId)
}
