package user

import (
	cloud "cloud.google.com/go/storage"
	"errors"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/responseerror"
	"github.com/Shresth92/audiophile/services"
	"github.com/Shresth92/audiophile/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"net/http"
	"strconv"
	"time"
)

type Controller struct {
	userService services.UserServices
}

func NewController(userService services.UserServices) *Controller {
	return &Controller{userService: userService}
}

func (c *Controller) AddProductToCart(ctx *gin.Context) {
	variantId := ctx.Param("variantId")
	userID := ctx.Value("userID").(string)
	err := c.userService.AddProductToCart(userID, variantId)
	if err != nil {
		logrus.Errorf("AddProductToCart: error in adding product to cart err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in adding product to cart")
		return
	}

	ctx.JSON(http.StatusOK, "Success")
}

func (c *Controller) UpdateProductCountInCart(ctx *gin.Context) {
	variantId := ctx.Param("variantId")
	count, err := strconv.ParseBool(ctx.Query("count"))
	if err != nil {
		logrus.Errorf("AddProductToCart: error in parsing count err: %v", err)
		responseerror.RespondClientErr(ctx, err, http.StatusBadRequest, "error in parsing count")
		return
	}

	userID := ctx.Value("userID").(string)
	err = c.userService.UpdateProductCountInCart(userID, variantId, count)
	if err != nil {
		logrus.Errorf("AddProductToCart: error in updating product count err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in updating product count")
		return
	}

	ctx.JSON(http.StatusOK, "Product count updated")
}

func (c *Controller) RemoveProductFromCart(ctx *gin.Context) {
	variantId := ctx.Param("variantId")
	userID := ctx.Value("userID").(string)
	err := c.userService.RemoveProductFromCart(userID, variantId)
	if err != nil {
		logrus.Errorf("RemoveProductFromCart: error in removing product from cart err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in removing product from cart")
		return
	}

	ctx.JSON(http.StatusOK, "Product removed from user")
}

func (c *Controller) DeleteMyCart(ctx *gin.Context) {
	userID := ctx.Value("userID").(string)
	err := c.userService.DeleteCart(userID)
	if err != nil {
		logrus.Errorf("DeleteMyCart: error in deleting my cart err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in deleting my cart")
		return
	}

	ctx.JSON(http.StatusOK, "user deletion successful")
}

func (c *Controller) GetMyCart(ctx *gin.Context) {
	userID := ctx.Value("userID").(string)
	cartItems, err := c.userService.GetCartProducts(userID)
	if err != nil {
		logrus.Errorf("GetMyCart: error in getting my cart err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in getting my cart")
		return
	}

	ctx.JSON(http.StatusOK, cartItems)
}

func (c *Controller) GetAllProducts(ctx *gin.Context) {
	limit, page, err := utils.GetLimitPage(ctx)
	if err != nil {
		logrus.Errorf("GetAllProducts: error in parsing limit and page err: %v", err)
		responseerror.RespondClientErr(ctx, err, http.StatusBadRequest, "error in parsing limit and page")
		return
	}

	searchString := ctx.Query("searchString")
	categoryFilter := ctx.Query("categoryFilter")
	brandFilter := ctx.Query("brandFilter")

	eg := &errgroup.Group{}
	var count int64
	var products []models.AllProducts

	eg.Go(func() error {
		count, err = c.userService.FilterAllProductsCount(searchString, categoryFilter, brandFilter)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		products, err = c.userService.FilterAllProducts(limit, page, searchString, categoryFilter, brandFilter)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		logrus.Errorf("GetAllProducts: error in getting products: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in getting products")
		return
	}

	productMap := make(map[string]models.Products)
	client := models.FirebaseClient
	for _, product := range products {
		_, ok := productMap[product.VariantID]

		signedUrl := &cloud.SignedURLOptions{
			Scheme:  cloud.SigningSchemeV4,
			Method:  "GET",
			Expires: time.Now().Add(15 * time.Minute),
		}
		url, err := client.Storage.Bucket(product.BucketName).SignedURL(product.Path, signedUrl)
		if err != nil {
			logrus.Errorf("GetAllProducts: error in generating image url: %v", err)
			responseerror.RespondGenericServerErr(ctx, err, "error in generating image")
			return
		}

		var newProduct models.Products
		if !ok {
			newProduct = models.Products{
				Id:           product.VariantID,
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
				Id:           product.VariantID,
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
				ImageLinks:   append(productMap[product.VariantID].ImageLinks, url),
			}
		}
		productMap[product.VariantID] = newProduct
	}

	ctx.JSON(http.StatusOK, models.Response{
		TotalRows: count,
		Rows:      productMap,
	})
}

func (c *Controller) GetProduct(ctx *gin.Context) {
	productID := ctx.Param("productId")
	products, err := c.userService.GetProduct(productID)
	if err != nil {
		logrus.Errorf("GetProduct: error in getting product err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in getting product")
		return
	}

	productMap := make(map[string]models.Products)
	client := models.FirebaseClient
	for _, product := range products {
		_, ok := productMap[product.VariantID]

		signedUrl := &cloud.SignedURLOptions{
			Scheme:  cloud.SigningSchemeV4,
			Method:  "GET",
			Expires: time.Now().Add(15 * time.Minute),
		}
		url, err := client.Storage.Bucket(product.BucketName).SignedURL(product.Path, signedUrl)
		if err != nil {
			logrus.Errorf("GetAllProducts: error in generating image url err: %v", err)
			responseerror.RespondGenericServerErr(ctx, err, "error in generating image url")
			return
		}

		if !ok {
			newProduct := models.Products{
				Id:           product.VariantID,
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
			productMap[product.VariantID] = newProduct
		} else {
			newProduct := models.Products{
				Id:           product.VariantID,
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
				ImageLinks:   append(productMap[product.VariantID].ImageLinks, url),
			}
			productMap[product.VariantID] = newProduct
		}
	}

	ctx.JSON(http.StatusOK, models.Response{
		TotalRows: int64(len(productMap)),
		Rows:      productMap},
	)
}

func (c *Controller) OrderProductByCart(ctx *gin.Context) {
	couponCode := ctx.Query("couponCode")
	addressID := ctx.Query("addressId")
	var variantIds []string
	userID := ctx.Value("userID").(string)
	cartItems, err := c.userService.GetCartProducts(userID)
	if err != nil {
		logrus.Errorf("OrderProductByCart: error in getting user products err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in getting user products")
		return
	}

	for _, cartItem := range cartItems {
		variantIds = append(variantIds, cartItem.VariantId)
		err := c.userService.UpdateProductStock(cartItem.VariantId, cartItem.Count)
		if err != nil {
			logrus.Errorf("OrderProductByCart: error in updating user products stock err: %v", err)
			responseerror.RespondGenericServerErr(ctx, err, "error in updating user products stock")
			return
		}
	}

	totalPrice, err := c.userService.GetTotalProductCost(variantIds)
	if err != nil {
		logrus.Errorf("OrderProductByCart: error in getting user products cost err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in getting user products cost")
		return
	}

	if couponCode != "" {
		totalPrice, err = c.userService.PriceAfterDiscount(totalPrice, couponCode)
		if err != nil {
			logrus.Errorf("OrderProductByCart: error in getting discounted price err: %v", err)
			responseerror.RespondGenericServerErr(ctx, err, "error in getting discounted price")
			return
		}
	}

	orderId, err := c.userService.GenerateOrderIdByCart(totalPrice, userID, addressID)
	if err != nil {
		logrus.Errorf("OrderProductByCart: error in generating order id err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in generating order id")
		return
	}

	var OrderedProducts []models.ProductOrdered
	for _, cartItem := range cartItems {
		productOrderId := uuid.New().String()
		OrderedProduct := models.ProductOrdered{
			ID:        productOrderId,
			VariantId: cartItem.VariantId,
			Quantity:  cartItem.Count,
			OrderId:   orderId,
		}
		OrderedProducts = append(OrderedProducts, OrderedProduct)
	}

	err = c.userService.AddProductsInOrder(OrderedProducts)
	if err != nil {
		logrus.Errorf("OrderProductByCart: error in adding products in order err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in adding products in order")
		return
	}

	err = c.userService.DeleteCart(userID)
	if err != nil {
		logrus.Errorf("OrderProductByCart: error in removing items from user err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in removing items from user")
		return
	}

	ctx.JSON(http.StatusOK, orderId)
}

func (c *Controller) GetAllOffers(ctx *gin.Context) {
	offers, err := c.userService.GetAllOffers()
	if err != nil {
		logrus.Errorf("GetAllOffers: error in getting all offers err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in getting all offers")
		return
	}

	ctx.JSON(http.StatusOK, offers)
}

func (c *Controller) GetMyOrders(ctx *gin.Context) {
	limit, page, err := utils.GetLimitPage(ctx)
	if err != nil {
		logrus.Errorf("GetMyFilteredProducts: error in parsing limit and page err: %v", err)
		responseerror.RespondClientErr(ctx, err, http.StatusBadRequest, "error in parsing limit and page")
		return
	}

	productStatusFilter := ctx.Query("productStatusFilter")
	userID := ctx.Value("userID").(string)

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
		responseerror.RespondClientErr(ctx, errors.New("there is no delivery status"), http.StatusBadRequest, "there is no delivery status")
		return
	}

	eg := &errgroup.Group{}
	var orders []models.Orders
	var ordersCount int64

	eg.Go(func() error {
		orders, err = c.userService.FilterMyOrders(userID, deliveryStatus, limit, page)
		if err != nil {
			return err
		}
		return nil
	})

	eg.Go(func() error {
		ordersCount, err = c.userService.CountFilterMyOrders(userID, deliveryStatus)
		if err != nil {
			return err
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		logrus.Errorf("GetMyFilteredProducts: error in getting my orders count err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in getting my orders count")
		return
	}

	ctx.JSON(http.StatusOK, models.Response{
		TotalRows: ordersCount,
		Rows:      orders,
	})
}

func (c *Controller) AddAddress(ctx *gin.Context) {
	addressDetails := models.Address{}
	if parseErr := ctx.ShouldBind(&addressDetails); parseErr != nil {
		responseerror.RespondClientErr(ctx, parseErr, http.StatusBadRequest, "error in parsing address")
		return
	}

	userID := ctx.Value("userID").(string)
	addressId, err := c.userService.AddAddress(userID, &addressDetails)
	if err != nil {
		logrus.Errorf("AddAddress: error in adding address err: %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in adding address")
		return
	}

	ctx.JSON(http.StatusOK, addressId)
}
