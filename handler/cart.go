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
	"net/http"
	"strconv"
)

func AddProductToCart(w http.ResponseWriter, r *http.Request) {
	variantId, err := uuid.Parse(chi.URLParam(r, "variantId"))
	if err != nil {
		logrus.Errorf("AddProductToCart: error in parsing variant id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing variant id")
		return
	}

	userCtx := r.Context().Value(middlewares.ContextUserKey).(*models.Claims)
	cartId, err := dbHelpers.AddProductToCart(userCtx.UserId, variantId, 1)
	if err != nil {
		logrus.Errorf("AddProductToCart: error in adding product to cart err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in adding product to cart")
		return
	}

	utils.Respond(w, http.StatusOK, cartId)
}

func UpdateProductCountInCart(w http.ResponseWriter, r *http.Request) {
	variantId, err := uuid.Parse(chi.URLParam(r, "variantId"))
	if err != nil {
		logrus.Errorf("AddProductToCart: error in parsing variant id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing variant id")
		return
	}

	count, err := strconv.ParseBool(r.URL.Query().Get("count"))
	if err != nil {
		logrus.Errorf("AddProductToCart: error in parsing count err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing count")
		return
	}

	userCtx := r.Context().Value(middlewares.ContextUserKey).(*models.Claims)
	err = dbHelpers.UpdateProductCountInCart(userCtx.UserId, variantId, count)
	if err != nil {
		logrus.Errorf("AddProductToCart: error in updating product count err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in updating product count")
		return
	}

	utils.Respond(w, http.StatusOK, "Product count updated")
}

func RemoveProductFromCart(w http.ResponseWriter, r *http.Request) {
	variantId, err := uuid.Parse(chi.URLParam(r, "variantId"))
	if err != nil {
		logrus.Errorf("RemoveProductFromCart: error in parsing variant id err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "error in parsing variant id")
		return
	}

	userCtx := r.Context().Value(middlewares.ContextUserKey).(*models.Claims)
	err = dbHelpers.RemoveCartProduct(userCtx.UserId, variantId)
	if err != nil {
		logrus.Errorf("RemoveProductFromCart: error in removing product from cart err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in removing product from cart")
		return
	}

	utils.Respond(w, http.StatusOK, "Product removed from cart")
}

func DeleteMyCart(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(middlewares.ContextUserKey).(*models.Claims)

	db := database.Db
	err := dbHelpers.DeleteCart(db, userCtx.UserId)
	if err != nil {
		logrus.Errorf("DeleteMyCart: error in deleting my cart err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "cart deletion failed")
		return
	}

	utils.Respond(w, http.StatusOK, "cart deletion successful")
}

func GetMyCart(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(middlewares.ContextUserKey).(*models.Claims)
	cart, err := dbHelpers.GetCartProducts(userCtx.UserId)
	if err != nil {
		logrus.Errorf("GetMyCart: error in getting my cart err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in getting my cart")
		return
	}

	utils.Respond(w, http.StatusOK, cart)
}
