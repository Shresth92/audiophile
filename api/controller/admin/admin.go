package admin

import (
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/responseerror"
	"github.com/Shresth92/audiophile/services"
	"github.com/Shresth92/audiophile/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"strconv"
	"time"
)

type Controller struct {
	adminService services.AdminServices
}

func NewController(adminService services.AdminServices) *Controller {
	return &Controller{adminService: adminService}
}

func (c *Controller) CreateProduct(ctx *gin.Context) {
	productDetails := models.ProductBody{}
	if parseErr := ctx.ShouldBind(&productDetails); parseErr != nil {
		responseerror.RespondClientErr(ctx, parseErr, http.StatusBadRequest, "error in parsing products")
		return
	}

	productID, productErr := c.adminService.CreateProduct(&productDetails)
	if productErr != nil {
		logrus.Errorf("CreateProduct: error in creating product err = %v", productErr)
		responseerror.RespondGenericServerErr(ctx, productErr, "error in creating products")
		_ = ctx.AbortWithError(http.StatusInternalServerError, productErr)
		return
	}

	variantId, variantErr := c.adminService.CreateVariant(productID, productDetails.Colour, productDetails.Stock, productDetails.Price)
	if variantErr != nil {
		logrus.Errorf("CreateProduct: error in creating variant err = %v", variantErr)
		responseerror.RespondGenericServerErr(ctx, variantErr, "error in creating variant")
		return
	}

	imageErr := c.adminService.UploadVariantImages(variantId, productDetails.ImageIds)
	if imageErr != nil {
		logrus.Errorf("CreateProduct: error in creating images err = %v", imageErr)
		responseerror.RespondGenericServerErr(ctx, imageErr, "error in uploading images")
		return
	}

	ctx.JSON(http.StatusCreated, "your product is successfully created")
}

func (c *Controller) CreateVariant(ctx *gin.Context) {
	variantDetails := models.VariantBody{}
	if parseErr := ctx.ShouldBind(&variantDetails); parseErr != nil {
		responseerror.RespondClientErr(ctx, parseErr, http.StatusBadRequest, "error in parsing variants")
		return
	}

	productID := ctx.Param("productId")

	variantId, variantErr := c.adminService.CreateVariant(productID, variantDetails.Colour, variantDetails.Stock, variantDetails.Price)
	if variantErr != nil {
		logrus.Errorf("CreateProduct: error in creating variant err = %v", variantErr)
		responseerror.RespondGenericServerErr(ctx, variantErr, "error in creating variant")
		return
	}

	uploadVariantImagesErr := c.adminService.UploadVariantImages(variantId, variantDetails.ImageIds)
	if uploadVariantImagesErr != nil {
		logrus.Errorf("CreateProduct: error in creating images err = %v", uploadVariantImagesErr)
		responseerror.RespondGenericServerErr(ctx, uploadVariantImagesErr, "error in uploading variant images")
		return
	}

	ctx.JSON(http.StatusCreated, "your utils is successfully created")
}

func (c *Controller) UploadImages(ctx *gin.Context) {
	client := models.FirebaseClient

	file, fileHeader, err := ctx.Request.FormFile("image")
	if err != nil {
		logrus.Errorf("UploadImages: error in parsing multipart form err = %v", err)
		responseerror.RespondClientErr(ctx, err, http.StatusBadRequest, "error in parsing image")
		return
	}

	err = ctx.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		logrus.Errorf("UploadImages: error in parsing multipart form err = %v", err)
		responseerror.RespondClientErr(ctx, err, http.StatusBadRequest, "error in parsing image")
		return
	}

	defer file.Close()
	imagePath := fileHeader.Filename + strconv.Itoa(int(time.Now().Unix()))
	bucket := utils.GetEnvValue("firebaseBucket")
	bucketStorage := client.Storage.Bucket(bucket).Object(imagePath).NewWriter(client.Ctx)

	_, err = io.Copy(bucketStorage, file)
	if err != nil {
		logrus.Errorf("UploadImages: error in file copying err: %v", err)
		responseerror.RespondClientErr(ctx, err, http.StatusBadRequest, "error in copying image")
		return
	}

	imageId, imageUploadErr := c.adminService.UploadImageFirebase(bucket, imagePath)
	if imageUploadErr != nil {
		logrus.Errorf("UploadImages: error in uploading image to firebase err = %v", imageUploadErr)
		responseerror.RespondGenericServerErr(ctx, imageUploadErr, "error in uploading images")
		return
	}

	if bucketErr := bucketStorage.Close(); err != nil {
		logrus.Errorf("UploadImages: error in closing firebase bucket err = %v", bucketErr)
		responseerror.RespondGenericServerErr(ctx, bucketErr, "error in closing firebase bucket")
		return
	}

	ctx.JSON(http.StatusCreated, imageId)
}

func (c *Controller) CreateCategory(ctx *gin.Context) {
	categoryDetails := models.Category{}
	if parseErr := ctx.ShouldBind(&categoryDetails); parseErr != nil {
		responseerror.RespondClientErr(ctx, parseErr, http.StatusBadRequest, "error in parsing categories")
		return
	}

	categoryErr := c.adminService.CreateCategory(categoryDetails.CategoryName)
	if categoryErr != nil {
		logrus.Errorf("CreateCategory: error in creating category err: %v", categoryErr)
		responseerror.RespondGenericServerErr(ctx, categoryErr, "error in creating category")
		return
	}

	ctx.JSON(http.StatusCreated, "Success")
}

func (c *Controller) CreateOffer(ctx *gin.Context) {
	offerDetails := models.Offer{}
	if parseErr := ctx.ShouldBind(&offerDetails); parseErr != nil {
		responseerror.RespondClientErr(ctx, parseErr, http.StatusBadRequest, "error in parsing offers")
		return
	}

	err := c.adminService.CreateOffer(&offerDetails)
	if err != nil {
		logrus.Errorf("CreateOffer: error in creating offer err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in creating offer")
		return
	}

	ctx.JSON(http.StatusCreated, "Success")
}

func (c *Controller) CreateBrand(ctx *gin.Context) {
	brandDetails := models.Brand{}
	if parseErr := ctx.ShouldBind(&brandDetails); parseErr != nil {
		responseerror.RespondClientErr(ctx, parseErr, http.StatusBadRequest, "error in parsing brands")
		return
	}

	err := c.adminService.CreateBrand(brandDetails.BrandName)
	if err != nil {
		logrus.Errorf("CreateBrand: error in creating brand err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in creating brand")
		return
	}

	ctx.JSON(http.StatusCreated, "Success")
}

func (c *Controller) DeleteVariant(ctx *gin.Context) {
	productID := ctx.Param("productID")
	variantID := ctx.Param("variantID")
	variantErr := c.adminService.DeleteVariant(productID, variantID)
	if variantErr != nil {
		logrus.Errorf("DeleteVariant: error in deleting variant err: %v", variantErr)
		responseerror.RespondGenericServerErr(ctx, variantErr, "error in deleting variant")
		return
	}

	ctx.JSON(http.StatusOK, "variant deleted successfully")
}

func (c *Controller) DeleteProduct(ctx *gin.Context) {
	productID := ctx.Param("productID")
	if err := c.adminService.DeleteProduct(productID); err != nil {
		logrus.Errorf("DeleteProduct: error in deleting product err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in deleting product")
		return
	}

	ctx.JSON(http.StatusOK, "variant deleted successfully")
}

func (c *Controller) DeleteCategory(ctx *gin.Context) {
	categoryId := ctx.Param("categoryId")
	if err := c.adminService.DeleteCategory(categoryId); err != nil {
		logrus.Errorf("DeleteCategory: error in deleting category err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in deleting category")
		return
	}

	ctx.JSON(http.StatusOK, "category deleted successfully")
}

func (c *Controller) DeleteBrand(ctx *gin.Context) {
	brandId := ctx.Param("brandId")
	if err := c.adminService.DeleteBrand(brandId); err != nil {
		logrus.Errorf("DeleteBrand: error in deleting brand err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in deleting brand")
		return
	}

	ctx.JSON(http.StatusOK, "brand deleted successfully")
}

func (c *Controller) UpdateProduct(ctx *gin.Context) {
	productID := ctx.Param("productId")
	productDetails := models.Product{}
	if parseErr := ctx.ShouldBind(&productDetails); parseErr != nil {
		responseerror.RespondClientErr(ctx, parseErr, http.StatusBadRequest, "error in parsing product details")
		return
	}

	if err := c.adminService.UpdateProduct(productID, &productDetails); err != nil {
		logrus.Errorf("UpdateProduct: error in updating product err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in updating product")
		return
	}

	ctx.JSON(http.StatusCreated, "product updated successfully")
}

func (c *Controller) UpdateVariant(ctx *gin.Context) {
	productID := ctx.Param("productId")
	variantID := ctx.Param("variantId")
	variantDetails := models.Variants{}
	if parseErr := ctx.ShouldBind(&variantDetails); parseErr != nil {
		responseerror.RespondClientErr(ctx, parseErr, http.StatusBadRequest, "error in parsing variant details")
		return
	}

	if err := c.adminService.UpdateVariant(productID, variantID, &variantDetails); err != nil {
		logrus.Errorf("UpdateVariant: error in updating variant err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in updating variant")
		return
	}

	ctx.JSON(http.StatusOK, "variant updated successfully")
}

func (c *Controller) UpdateCategory(ctx *gin.Context) {
	categoryID := ctx.Param("categoryId")
	categoryDetails := models.Category{}
	if parseErr := ctx.ShouldBind(&categoryDetails); parseErr != nil {
		responseerror.RespondClientErr(ctx, parseErr, http.StatusBadRequest, "error in parsing categories")
		return
	}

	if err := c.adminService.UpdateCategory(categoryID, categoryDetails.CategoryName); err != nil {
		logrus.Errorf("UpdateCategory: error in updating category err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in updating category")
	}

	ctx.JSON(http.StatusOK, "category updated successfully")
}

func (c *Controller) UpdateBrand(ctx *gin.Context) {
	brandID := ctx.Param("brandId")
	brandDetails := models.Brand{}
	if parseErr := ctx.ShouldBind(&brandDetails); parseErr != nil {
		responseerror.RespondClientErr(ctx, parseErr, http.StatusBadRequest, "error in parsing brands")
		return
	}

	if err := c.adminService.UpdateBrand(brandID, brandDetails.BrandName); err != nil {
		responseerror.RespondGenericServerErr(ctx, err, "error in updating brand")
		return
	}

	ctx.JSON(http.StatusOK, "brand updated successfully")
}

func (c *Controller) GetAllUsers(ctx *gin.Context) {
	limit, page, err := utils.GetLimitPage(ctx)
	if err != nil {
		logrus.Errorf("GetAllUsers: error in parsing limit and page err: %v", err)
		responseerror.RespondClientErr(ctx, err, http.StatusBadRequest, "error in parsing limit and page")
		return
	}

	eg := &errgroup.Group{}
	var users []models.Users
	var usersCount int64

	eg.Go(func() error {
		users, err = c.adminService.GetAllUsers(limit, page)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		usersCount, err = c.adminService.UsersCount()
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		logrus.Errorf("GetAllUsers: error in getting users err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in getting users")
		return
	}

	ctx.JSON(http.StatusOK, struct {
		TotalRows int64
		Rows      []models.Users
	}{
		TotalRows: usersCount,
		Rows:      users,
	})
}

func (c *Controller) GetAllBrands(ctx *gin.Context) {
	limit, page, err := utils.GetLimitPage(ctx)
	if err != nil {
		logrus.Errorf("GetAllBrands: error in parsing limit and page err: %v", err)
		responseerror.RespondClientErr(ctx, err, http.StatusBadRequest, "error in parsing limit and page")
		return
	}

	eg := &errgroup.Group{}
	var brands []models.Brand
	var brandsCount int64

	eg.Go(func() error {
		brands, err = c.adminService.GetAllBrands(limit, page)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		brandsCount, err = c.adminService.GetBrandsCount()
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		logrus.Errorf("GetAllBrands: error in getting all brands err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in getting brands")
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		TotalRows: brandsCount,
		Rows:      brands,
	})
}

func (c *Controller) GetAllCategory(ctx *gin.Context) {
	limit, page, err := utils.GetLimitPage(ctx)
	if err != nil {
		logrus.Errorf("GetAllCategory: error in parsing limit and page err: %v", err)
		responseerror.RespondClientErr(ctx, err, http.StatusBadRequest, "error in parsing limit and page")
		return
	}

	eg := &errgroup.Group{}
	var categories []models.Category
	var categoriesCount int64

	eg.Go(func() error {
		categories, err = c.adminService.GetAllCategory(limit, page)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		categoriesCount, err = c.adminService.GetCategoryCount()
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		logrus.Errorf("GetAllCategory: error in getting all categories err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in getting categories")
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		TotalRows: categoriesCount,
		Rows:      categories,
	})
}

func (c *Controller) ChangeUserRole(ctx *gin.Context) {
	userID := ctx.Param("userId")
	adminID := ctx.Value("userID").(string)
	if err := c.adminService.ChangeUserRole(userID, adminID); err != nil {
		logrus.Errorf("ChangeUserRole: error in changing user role err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in changing role")
		return
	}

	ctx.JSON(http.StatusOK, "user role changed")
}
