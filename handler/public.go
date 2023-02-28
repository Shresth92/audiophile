package handler

import (
	"errors"
	"github.com/Shresth92/audiophile/database"
	"github.com/Shresth92/audiophile/database/dbHelpers"
	"github.com/Shresth92/audiophile/middlewares"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/utils"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func Register(w http.ResponseWriter, r *http.Request) {
	body := models.Users{}
	if err := utils.DecodeBody(r.Body, &body); err != nil {
		logrus.Errorf("Register: error in parsing body err = %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "Body parse failed")
		return
	}

	hashedPassword, err := utils.HashPassword(body.Password)
	if err != nil {
		logrus.Errorf("Register: error in password hashing err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in password hashing")
		return
	}

	isUserExists, err := dbHelpers.CheckUserExist(body.Email)
	if err != nil {
		logrus.Errorf("Register: error in getting user credentials err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in getting user credentials")
		return
	}

	if isUserExists {
		logrus.Errorf("Register: user already exists err: %v", err)
		utils.RespondError(w, http.StatusBadRequest, errors.New("user already exists"), false, "user already exists")
		return
	}

	tx := database.Db.Begin()
	userId, err := dbHelpers.CreateUser(tx, body.Email, hashedPassword)
	if err != nil {
		logrus.Errorf("Register: error in creating user err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in creating user")
		tx.Rollback()
		return
	}

	_, err = dbHelpers.CreateUserRole(tx, userId, models.User)
	if err != nil {
		logrus.Errorf("Register: error in creating role err: %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in creating role")
		tx.Rollback()
		return
	}

	tx.Commit()
	utils.Respond(w, http.StatusOK, userId)
}

func Login(w http.ResponseWriter, r *http.Request) {
	body := models.Users{}
	if err := utils.DecodeBody(r.Body, &body); err != nil {
		logrus.Errorf("Register: error in parsing body err = %v", err)
		utils.RespondError(w, http.StatusBadRequest, err, true, "Body parse failed")
		return
	}

	var role models.Roles

	if strings.Contains(r.URL.String(), "admin-login") {
		role = models.Admin
	} else {
		role = models.User
	}

	user, err := dbHelpers.GetUserDetails(body.Email)
	if err != nil {
		logrus.Errorf("Login: error in getting user credentials err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in getting user credentials")
		return
	}

	isPasswordMatched := utils.CheckPassword(body.Password, user.Password)
	if !isPasswordMatched {
		logrus.Errorf("Login: error in password matching err = %v", err)
		utils.RespondError(w, http.StatusBadRequest, errors.New("password not matched"), true, "error in password matching")
		return
	}

	sessionId, err := dbHelpers.GetSessionId(user.Id)
	if err != nil {
		logrus.Errorf("Login: error in creating session err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in creating session")
		return
	}

	token, err := utils.GenerateJWTToken(user.Id, sessionId, role)
	if err != nil {
		logrus.Errorf("Login: error in generating jwt token err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in generating jwt token")
		return
	}

	utils.Respond(w, http.StatusOK, token)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	userCtx := r.Context().Value(middlewares.ContextUserKey).(*models.Claims)
	err := dbHelpers.LogOut(userCtx.UserId, userCtx.SessionId)
	if err != nil {
		logrus.Errorf("Logout: error in logging out err = %v", err)
		utils.RespondError(w, http.StatusInternalServerError, err, false, "error in logging out")
		return
	}

	utils.Respond(w, http.StatusOK, "You are logged out")
}
