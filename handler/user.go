package handler

import (
	cloud "cloud.google.com/go/storage"
	"errors"
	"github.com/Shresth92/audiophile/database"
	"github.com/Shresth92/audiophile/database/dbHelpers"
	"github.com/Shresth92/audiophile/middlewares"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/utils"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"time"
)

func GetAllProducts(w http.ResponseWriter, r *http.Request) {
	limit, page, err := utils.GetLimitPage(r.URL.Query())
	if err != nil {
		logrus.Errorf("GetAllProducts: error in parsing limit and page err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing page")
		return
	}

	searchString := r.URL.Query().Get("searchString")
	categoryFilter := r.URL.Query().Get("categoryFilter")
	brandFilter := r.URL.Query().Get("brandFilter")

	eg := &errgroup.Group{}
	var count int64
	var products []models.AllProducts

	eg.Go(func() error {
		count, err = dbHelpers.FilterAllProductsCount(searchString, categoryFilter, brandFilter)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		products, err = dbHelpers.FilterAllProducts(limit, page, searchString, categoryFilter, brandFilter)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		logrus.Errorf("GetAllProducts: error in getting products: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in getting products")
		return
	}

	productMap := make(map[uuid.UUID]models.Products)
	client := models.FirebaseClient
	for _, product := range products {
		_, ok := productMap[product.Id]

		signedUrl := &cloud.SignedURLOptions{
			Scheme:  cloud.SigningSchemeV4,
			Method:  "GET",
			Expires: time.Now().Add(15 * time.Minute),
		}
		url, err := client.Storage.Bucket(product.BucketName).SignedURL(product.Path, signedUrl)
		if err != nil {
			logrus.Errorf("GetAllProducts: error in generating image url: %v", err)
			utils.RespondError(w, http.StatusInternalServerError, err, false, "error in generating image url")
			return
		}

		var newProduct models.Products
		if !ok {
			newProduct = models.Products{
				Id:           product.Id,
				ProductName:  product.ProductName,
				ModelName:    product.ModelName,
				BrandName:    product.BrandName,
				CategoryName: product.CategoryName,
				Return:       product.Return,
				Warranty:     product.Warranty,
				Wireless:     product.Wireless,
				Colour:       product.Colour,
				Price:        product.Price,
				Stock:        product.Stock,
				ImageLinks:   []string{url},
			}
		} else {
			newProduct = models.Products{
				Id:           product.Id,
				ProductName:  product.ProductName,
				ModelName:    product.ModelName,
				BrandName:    product.BrandName,
				CategoryName: product.CategoryName,
				Return:       product.Return,
				Warranty:     product.Warranty,
				Wireless:     product.Wireless,
				Colour:       product.Colour,
				Price:        product.Price,
				Stock:        product.Stock,
				ImageLinks:   append(productMap[product.Id].ImageLinks, url),
			}
		}
		productMap[product.Id] = newProduct
	}

	utils.Respond(w, http.StatusOK, models.Response{
		TotalRows: count,
		Rows:      productMap,
	})
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := uuid.Parse(chi.URLParam(r, "productId"))
	if err != nil {
		logrus.Errorf("GetProduct: error in parsing product id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing product id")
		return
	}

	products, err := dbHelpers.GetProduct(productId)
	if err != nil {
		logrus.Errorf("GetProduct: error in getting product err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in getting product")
		return
	}

	productMap := make(map[uuid.UUID]models.Products)
	client := models.FirebaseClient
	for _, product := range products {
		_, ok := productMap[product.Id]

		signedUrl := &cloud.SignedURLOptions{
			Scheme:  cloud.SigningSchemeV4,
			Method:  "GET",
			Expires: time.Now().Add(15 * time.Minute),
		}
		url, err := client.Storage.Bucket(product.BucketName).SignedURL(product.Path, signedUrl)
		if err != nil {
			logrus.Errorf("GetAllProducts: error in generating image url err: %v", err)
			utils.RespondError(w, http.StatusInternalServerError, err, false, "error in generating image url")
			return
		}

		if !ok {
			newProduct := models.Products{
				Id:           product.Id,
				ProductName:  product.ProductName,
				ModelName:    product.ModelName,
				BrandName:    product.BrandName,
				CategoryName: product.CategoryName,
				Return:       product.Return,
				Warranty:     product.Warranty,
				Wireless:     product.Wireless,
				Colour:       product.Colour,
				Price:        product.Price,
				Stock:        product.Stock,
				ImageLinks:   []string{url},
			}
			productMap[product.Id] = newProduct
		} else {
			newProduct := models.Products{
				Id:           product.Id,
				ProductName:  product.ProductName,
				ModelName:    product.ModelName,
				BrandName:    product.BrandName,
				CategoryName: product.CategoryName,
				Return:       product.Return,
				Warranty:     product.Warranty,
				Wireless:     product.Wireless,
				Colour:       product.Colour,
				Price:        product.Price,
				Stock:        product.Stock,
				ImageLinks:   append(productMap[product.Id].ImageLinks, url),
			}
			productMap[product.Id] = newProduct
		}
	}

	utils.Respond(w, http.StatusOK, models.Response{
		TotalRows: int64(len(productMap)),
		Rows:      productMap},
	)
}

func OrderProductByCart(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	couponCode := query.Get("couponCode")
	addressId, err := uuid.Parse(query.Get("addressId"))
	if err != nil {
		logrus.Errorf("OrderProductByCart: error in parsing address id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "Address parsing failed")
		return
	}

	var variantIds []uuid.UUID
	userCtx := r.Context().Value(middlewares.ContextUserKey).(*models.Claims)

	cartItems, err := dbHelpers.GetCartProducts(userCtx.UserId)
	if err != nil {
		logrus.Errorf("OrderProductByCart: error in getting cart products err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in getting cart products")
		return
	}

	tx := database.Db.Begin()
	for _, cartItem := range cartItems {
		variantIds = append(variantIds, cartItem.VariantId)
		err := dbHelpers.UpdateProductStock(tx, cartItem.VariantId, cartItem.Count)
		if err != nil {
			logrus.Errorf("OrderProductByCart: error in updating cart products stock err: %v", err)
			utils.RespondError(w, http.StatusBadRequest, err, true, "error in updating cart products stock")
			tx.Rollback()
			return
		}
	}

	totalPrice, err := dbHelpers.GetTotalProductCost(variantIds)
	if err != nil {
		logrus.Errorf("OrderProductByCart: error in getting cart products cost err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in getting cart products cost")
		tx.Rollback()
		return
	}

	if couponCode != "" {
		totalPrice, err = dbHelpers.PriceAfterDiscount(totalPrice, couponCode)
		if err != nil {
			logrus.Errorf("OrderProductByCart: error in getting discounted price err: %v", err)
			utils.RespondError(w, http.StatusInternalServerError, err, false, "error in getting discounted price")
			tx.Rollback()
			return
		}
	}

	orderId, err := dbHelpers.GenerateOrderIdByCart(tx, totalPrice, userCtx.UserId, addressId)
	if err != nil {
		logrus.Errorf("OrderProductByCart: error in generating order id err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in generating order id")
		tx.Rollback()
		return
	}

	var OrderedProducts []models.ProductOrdered
	for _, cartItem := range cartItems {
		productOrderId := uuid.New()
		OrderedProduct := models.ProductOrdered{
			Id:        productOrderId,
			VariantId: cartItem.VariantId,
			Quantity:  cartItem.Count,
			OrderId:   orderId,
		}
		OrderedProducts = append(OrderedProducts, OrderedProduct)
	}

	err = dbHelpers.AddProductsInOrder(tx, OrderedProducts)
	if err != nil {
		logrus.Errorf("OrderProductByCart: error in adding products in order err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in adding products in order")
		tx.Rollback()
		return
	}

	err = dbHelpers.DeleteCart(tx, userCtx.UserId)
	if err != nil {
		logrus.Errorf("OrderProductByCart: error in removing items from cart err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in removing items from cart")
		tx.Rollback()
		return
	}

	tx.Commit()
	utils.Respond(w, http.StatusOK, orderId)
}

func GetAllOffers(w http.ResponseWriter, r *http.Request) {
	offers, err := dbHelpers.GetAllOffers()
	if err != nil {
		logrus.Errorf("GetAllOffers: error in getting all offers err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in getting all offers")
		return
	}

	utils.Respond(w, http.StatusOK, offers)
}

func GetMyOrders(w http.ResponseWriter, r *http.Request) {
	limit, page, err := utils.GetLimitPage(r.URL.Query())
	if err != nil {
		logrus.Errorf("GetMyFilteredProducts: error in parsing limit and page err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing page")
		return
	}

	productStatusFilter := r.URL.Query().Get("productStatusFilter")
	userCtx := r.Context().Value(middlewares.ContextUserKey).(*models.Claims)

	var deliveryStatus models.DeliveryStatus
	if productStatusFilter == "onTheWay" {
		deliveryStatus = models.OnTheWay
	} else if productStatusFilter == "delivered" {
		deliveryStatus = models.Delivered
	} else if productStatusFilter == "canceled" {
		deliveryStatus = models.Canceled
	} else if productStatusFilter == "return" {
		deliveryStatus = models.Return
	} else {
		logrus.Error("GetMyFilteredProducts: there is no delivery status")
		utils.RespondError(w, http.StatusBadRequest, errors.New("there is no delivery status"), true, "there is no delivery status")
		return
	}

	eg := &errgroup.Group{}
	var orders []models.Orders
	var ordersCount int64

	eg.Go(func() error {
		orders, err = dbHelpers.FilterMyOrders(userCtx.UserId, deliveryStatus, limit, page)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		ordersCount, err = dbHelpers.CountFilterMyOrders(userCtx.UserId, deliveryStatus)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		logrus.Errorf("GetMyFilteredProducts: error in getting my orders count err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in getting my orders count")
		return
	}

	utils.Respond(w, http.StatusOK, models.Response{
		TotalRows: ordersCount,
		Rows:      orders,
	})
}

func AddAddress(w http.ResponseWriter, r *http.Request) {
	body := models.Address{}
	if err := utils.DecodeBody(r.Body, &body); err != nil {
		logrus.Errorf("AddAddress: error in parsing body err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing body")
		return
	}

	userCtx := r.Context().Value(middlewares.ContextUserKey).(*models.Claims)
	addressId, err := dbHelpers.AddAddress(userCtx.UserId, &body)
	if err != nil {
		logrus.Errorf("AddAddress: error in adding address err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in adding address")
		return
	}

	utils.Respond(w, http.StatusOK, addressId)
}
