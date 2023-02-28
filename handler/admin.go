package handler

import (
	"github.com/Shresth92/audiophile/database"
	"github.com/Shresth92/audiophile/database/dbHelpers"
	"github.com/Shresth92/audiophile/middlewares"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io"
	"net/http"
	"strconv"
	"time"
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	body := models.ProductBody{}
	if err := utils.DecodeBody(r.Body, &body); err != nil {
		logrus.Errorf("CreateProduct: error in parsing body err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing body")
		return
	}

	tx := database.Db.Begin()
	productId, err := dbHelpers.CreateProduct(tx, &body)
	if err != nil {
		logrus.Errorf("CreateProduct: error in creating product err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in creating product")
		tx.Rollback()
		return
	}

	variantId, err := dbHelpers.CreateVariant(tx, productId, body.Colour, body.Stock, body.Price)
	if err != nil {
		logrus.Errorf("CreateProduct: error in creating variant err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in creating variant")
		tx.Rollback()
		return
	}

	_, err = dbHelpers.UploadVariantImages(tx, variantId, body.ImageIds)
	if err != nil {
		logrus.Errorf("CreateProduct: error in creating images err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in creating images")
		tx.Rollback()
		return
	}

	tx.Commit()
	utils.Respond(w, http.StatusCreated, "your product is successfully created")
}

func CreateVariant(w http.ResponseWriter, r *http.Request) {
	body := models.VariantBody{}
	if err := utils.DecodeBody(r.Body, &body); err != nil {
		logrus.Errorf("CreateProduct: error in parsing body err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing body")
		return
	}

	productId, err := uuid.Parse(chi.URLParam(r, "productId"))
	if err != nil {
		logrus.Errorf("DeleteProduct: error in parsing product id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing product id")
		return
	}

	tx := database.Db.Begin()

	variantId, err := dbHelpers.CreateVariant(tx, productId, body.Colour, body.Stock, body.Price)
	if err != nil {
		logrus.Errorf("CreateProduct: error in creating variant err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in creating variant")
		tx.Rollback()
		return
	}

	_, err = dbHelpers.UploadVariantImages(tx, variantId, body.ImageIds)
	if err != nil {
		logrus.Errorf("CreateProduct: error in creating images err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in creating images")
		tx.Rollback()
		return
	}

	tx.Commit()
	utils.Respond(w, http.StatusCreated, "your product is successfully created")
}

func UploadImages(w http.ResponseWriter, r *http.Request) {
	client := models.FirebaseClient

	file, fileHeader, err := r.FormFile("image")
	err = r.ParseMultipartForm(10 << 20)
	if err != nil {
		logrus.Errorf("UploadImages: error in parsing multipart form err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in parsing multipart form")
		return
	}

	defer file.Close()
	imagePath := fileHeader.Filename + strconv.Itoa(int(time.Now().Unix()))
	bucket := utils.GetEnvValue("firebaseBucket")
	bucketStorage := client.Storage.Bucket(bucket).Object(imagePath).NewWriter(client.Ctx)

	_, err = io.Copy(bucketStorage, file)
	if err != nil {
		logrus.Errorf("UploadImages: error in file copying err: %v", err)
		utils.RespondError(w, http.StatusBadGateway, err, true, "error in file copying err")
		return
	}

	imageId, err := dbHelpers.UploadImageFirebase(bucket, imagePath)
	if err != nil {
		logrus.Errorf("UploadImages: error in uploading image to firebase err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in uploading image to firebase")
		return
	}

	if err := bucketStorage.Close(); err != nil {
		logrus.Errorf("UploadImages: error in closing firebase bucket err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in closing firebase bucket")
		return
	}

	utils.Respond(w, http.StatusCreated, imageId)
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	body := models.Category{}
	if err := utils.DecodeBody(r.Body, &body); err != nil {
		logrus.Errorf("CreateCategory: error in parsing body err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing body")
		return
	}

	categoryId, err := dbHelpers.CreateCategory(body.CategoryName)
	if err != nil {
		logrus.Errorf("CreateCategory: error in creating category err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in creating category")
		return
	}

	utils.Respond(w, http.StatusCreated, categoryId)
}

func CreateOffer(w http.ResponseWriter, r *http.Request) {
	body := models.Offer{}
	if err := utils.DecodeBody(r.Body, &body); err != nil {
		logrus.Errorf("CreateOffer: error in parsing body err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing body")
		return
	}

	offerId, err := dbHelpers.CreateOffer(&body)
	if err != nil {
		logrus.Errorf("CreateOffer: error in creating offer err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in creating offer")
		return
	}

	utils.Respond(w, http.StatusCreated, offerId)
}

func CreateBrand(w http.ResponseWriter, r *http.Request) {
	body := models.Brand{}
	if err := utils.DecodeBody(r.Body, &body); err != nil {
		logrus.Errorf("CreateBrand: error in parsing body err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing body")
		return
	}

	brandId, err := dbHelpers.CreateBrand(body.BrandName)
	if err != nil {
		logrus.Errorf("CreateBrand: error in creating brand err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in creating brand")
		return
	}

	utils.Respond(w, http.StatusCreated, brandId)
}

func DeleteVariant(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		logrus.Errorf("DeleteVariant: error in parsing product id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing product id")
		return
	}

	variantID, err := uuid.Parse(chi.URLParam(r, "variantID"))
	if err != nil {
		logrus.Errorf("DeleteVariant: error in parsing variant id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing variant id")
		return
	}

	err = dbHelpers.DeleteVariant(productID, variantID)
	if err != nil {
		logrus.Errorf("DeleteVariant: error in deleting variant err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in deleting variant")
		return
	}

	utils.Respond(w, http.StatusOK, "variant deleted successfully")
}

func DeleteProduct(w http.ResponseWriter, r *http.Request) {
	productID, err := uuid.Parse(chi.URLParam(r, "productID"))
	if err != nil {
		logrus.Errorf("DeleteProduct: error in parsing product id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing product id")
		return
	}

	err = dbHelpers.DeleteProduct(productID)
	if err != nil {
		logrus.Errorf("DeleteProduct: error in deleting product err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in deleting product")
		return
	}

	utils.Respond(w, http.StatusOK, "variant deleted successfully")
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryId, err := uuid.Parse(chi.URLParam(r, "categoryId"))
	if err != nil {
		logrus.Errorf("DeleteCategory: error in parsing category id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing category id")
		return
	}

	err = dbHelpers.DeleteCategory(categoryId)
	if err != nil {
		logrus.Errorf("DeleteCategory: error in deleting category err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in deleting category")
		return
	}

	utils.Respond(w, http.StatusOK, "category deleted successfully")
}

func DeleteBrand(w http.ResponseWriter, r *http.Request) {
	brandId, err := uuid.Parse(chi.URLParam(r, "brandId"))
	if err != nil {
		logrus.Errorf("DeleteBrand: error in parsing brand id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing brand id")
		return
	}

	err = dbHelpers.DeleteBrand(brandId)
	if err != nil {
		logrus.Errorf("DeleteBrand: error in deleting brand err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in deleting brand")
		return
	}

	utils.Respond(w, http.StatusOK, "brand deleted successfully")
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := uuid.Parse(chi.URLParam(r, "productId"))
	if err != nil {
		logrus.Errorf("UpdateProduct: error in parsing product id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing product id")
		return
	}

	body := models.Product{}
	if err := utils.DecodeBody(r.Body, &body); err != nil {
		logrus.Errorf("UpdateProduct: error in parsing body err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing body")
		return
	}

	err = dbHelpers.UpdateProduct(productId, body.ProductName, body.ModelName, body.Return, body.Warranty, body.Wireless)
	if err != nil {
		logrus.Errorf("UpdateProduct: error in updating product err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in updating product")
		return
	}

	utils.Respond(w, http.StatusOK, "product updated successfully")
}

func UpdateVariant(w http.ResponseWriter, r *http.Request) {
	productId, err := uuid.Parse(chi.URLParam(r, "productId"))
	if err != nil {
		logrus.Errorf("UpdateVariant: error in parsing product id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing product id")
		return
	}

	variantId, err := uuid.Parse(chi.URLParam(r, "variantId"))
	if err != nil {
		logrus.Errorf("UpdateVariant: error in parsing variant id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing variant id")
		return
	}

	body := models.Variants{}
	if err := utils.DecodeBody(r.Body, &body); err != nil {
		logrus.Errorf("UpdateVariant: error in parsing body err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing body")
		return
	}

	err = dbHelpers.UpdateVariant(productId, variantId, body.Colour, body.Price, body.Stock)
	if err != nil {
		logrus.Errorf("UpdateVariant: error in updating variant err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in updating variant")
		return
	}

	utils.Respond(w, http.StatusOK, "variant updated successfully")
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	categoryId, err := uuid.Parse(chi.URLParam(r, "categoryId"))
	if err != nil {
		logrus.Errorf("UpdateCategory: error in parsing category id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing category id")
		return
	}

	body := models.Category{}
	if err := utils.DecodeBody(r.Body, &body); err != nil {
		logrus.Errorf("UpdateCategory: error in parsing body err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "UpdateCategory: error in parsing body err")
		return
	}

	err = dbHelpers.UpdateCategory(categoryId, body.CategoryName)
	if err != nil {
		logrus.Errorf("UpdateCategory: error in updating category err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in updating category")
		return
	}

	utils.Respond(w, http.StatusOK, "category updated successfully")
}

func UpdateBrand(w http.ResponseWriter, r *http.Request) {
	brandId, err := uuid.Parse(chi.URLParam(r, "brandId"))
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, true, "brand id parsing failed")
		return
	}

	body := struct {
		BrandName string `json:"brandName"`
	}{}
	if err := utils.DecodeBody(r.Body, &body); err != nil {
		utils.RespondError(w, http.StatusBadRequest, err, true, "Body parse failed")
		return
	}

	err = dbHelpers.UpdateBrand(brandId, body.BrandName)
	if err != nil {
		utils.RespondError(w, http.StatusInternalServerError, err, false, "brand update failed")
		return
	}
	utils.Respond(w, http.StatusOK, "brand updated successfully")
}

func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	limit, page, err := utils.GetLimitPage(r.URL.Query())
	if err != nil {
		logrus.Errorf("GetAllUsers: error in parsing limit and page err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing page")
		return
	}

	eg := &errgroup.Group{}
	var users []models.Users
	var usersCount int64

	eg.Go(func() error {
		users, err = dbHelpers.GetAllUsers(limit, page)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		usersCount, err = dbHelpers.UsersCount()
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		logrus.Errorf("GetAllUsers: error in getting users err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in getting users")
		return
	}

	utils.Respond(w, http.StatusOK, struct {
		TotalRows int64
		Rows      []models.Users
	}{
		TotalRows: usersCount,
		Rows:      users,
	})
}

func GetAllBrands(w http.ResponseWriter, r *http.Request) {
	limit, page, err := utils.GetLimitPage(r.URL.Query())
	if err != nil {
		logrus.Errorf("GetAllBrands: error in parsing limit and page err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing page")
		return
	}

	eg := &errgroup.Group{}
	var brands []models.Brand
	var brandsCount int64

	eg.Go(func() error {
		brands, err = dbHelpers.GetAllBrands(limit, page)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		brandsCount, err = dbHelpers.GetBrandsCount()
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		logrus.Errorf("GetAllBrands: error in getting all brands err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in getting all brands")
		return
	}

	utils.Respond(w, http.StatusOK, models.Response{
		TotalRows: brandsCount,
		Rows:      brands,
	})
}

func GetAllCategory(w http.ResponseWriter, r *http.Request) {
	limit, page, err := utils.GetLimitPage(r.URL.Query())
	if err != nil {
		logrus.Errorf("GetAllCategory: error in parsing limit and page err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing page")
		return
	}

	eg := &errgroup.Group{}
	var categories []models.Category
	var categoriesCount int64

	eg.Go(func() error {
		categories, err = dbHelpers.GetAllCategory(limit, page)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		categoriesCount, err = dbHelpers.GetCategoryCount()
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		logrus.Errorf("GetAllCategory: error in getting all categories err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in getting all categories")
		return
	}

	utils.Respond(w, http.StatusOK, models.Response{
		TotalRows: categoriesCount,
		Rows:      categories,
	})
}

func ChangeUserRole(w http.ResponseWriter, r *http.Request) {
	userId, err := uuid.Parse(chi.URLParam(r, "userId"))
	if err != nil {
		logrus.Errorf("ChangeUserRole: error in parsing user id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing user id")
		return
	}

	userCtx := r.Context().Value(middlewares.ContextUserKey).(*models.Claims)
	err = dbHelpers.ChangeUserRole(userId, userCtx.UserId)
	if err != nil {
		logrus.Errorf("ChangeUserRole: error in changing user role err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in changing user role")
		return
	}

	utils.Respond(w, http.StatusOK, "user role changed")
}
