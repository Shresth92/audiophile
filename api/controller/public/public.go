package public

import (
	"errors"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/responseerror"
	"github.com/Shresth92/audiophile/services"
	"github.com/Shresth92/audiophile/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

type Controller struct {
	publicService services.PublicService
}

func NewController(publicService services.PublicService) *Controller {
	return &Controller{publicService: publicService}
}

func (c *Controller) Register(ctx *gin.Context) {
	userDetails := models.Users{}
	if parseErr := ctx.ShouldBind(&userDetails); parseErr != nil {
		responseerror.RespondClientErr(ctx, parseErr, http.StatusBadRequest, "error in parsing user details")
		return
	}

	hashedPassword, hashedPasswordErr := utils.HashPassword(userDetails.Password)
	if hashedPasswordErr != nil {
		logrus.Errorf("Register: error in password hashing err = %v", hashedPasswordErr)
		responseerror.RespondGenericServerErr(ctx, errors.New("error in password hashing"), "error in password hashing")
		return
	}

	isUserExists, isUserExistsErr := c.publicService.CheckUserExist(userDetails.Email)
	if isUserExistsErr != nil {
		logrus.Errorf("Register: error in getting user credentials err = %v", isUserExistsErr)
		responseerror.RespondGenericServerErr(ctx, isUserExistsErr, "error in getting user credentials")
		return
	}

	if isUserExists {
		responseerror.RespondClientErr(ctx, errors.New("user already exists"), http.StatusConflict, "user already exists")
		_ = ctx.AbortWithError(http.StatusConflict, errors.New("user already exists"))
		return
	}

	userId, userErr := c.publicService.CreateUser(userDetails.Email, hashedPassword)
	if userErr != nil {
		logrus.Errorf("Register: error in creating user err: %v", userErr)
		responseerror.RespondGenericServerErr(ctx, userErr, "error in creating user")
		_ = ctx.AbortWithError(http.StatusInternalServerError, userErr)
		return
	}

	_, roleErr := c.publicService.CreateUserRole(userId, models.User)
	if roleErr != nil {
		logrus.Errorf("Register: error in creating role err: %v", roleErr)
		responseerror.RespondGenericServerErr(ctx, roleErr, "error in creating role")
		_ = ctx.AbortWithError(http.StatusInternalServerError, roleErr)
		return
	}

	ctx.JSON(http.StatusCreated, userId)
}

func (c *Controller) Login(ctx *gin.Context) {
	userDetails := models.Users{}
	if parseErr := ctx.ShouldBind(&userDetails); parseErr != nil {
		responseerror.RespondClientErr(ctx, parseErr, http.StatusBadRequest, "error in parsing user details")
		return
	}

	var role models.Roles

	if strings.Contains(ctx.Request.URL.String(), "admin-login") {
		role = models.Admin
	} else {
		role = models.User
	}

	user, userErr := c.publicService.GetUserDetails(userDetails.Email)
	if userErr != nil {
		logrus.Errorf("Login: error in getting user credentials err = %v", userErr)
		responseerror.RespondGenericServerErr(ctx, userErr, "error in getting user details")
		return
	}

	isPasswordMatched := utils.CheckPassword(userDetails.Password, user.Password)
	if !isPasswordMatched {
		responseerror.RespondGenericServerErr(ctx, userErr, "incorrect password")
		return
	}

	sessionId, sessionErr := c.publicService.GetSessionId(user.Id)
	if sessionErr != nil {
		logrus.Errorf("Login: error in creating session err = %v", sessionErr)
		responseerror.RespondGenericServerErr(ctx, sessionErr, "error in creating session")
		return
	}

	token, tokenErr := utils.GenerateJWTToken(user.Id, sessionId, role)
	if tokenErr != nil {
		logrus.Errorf("Login: error in generating jwt token err = %v", tokenErr)
		responseerror.RespondGenericServerErr(ctx, tokenErr, "error in generating jwt token")
		return
	}

	ctx.JSON(http.StatusCreated, token)
}

func (c *Controller) Logout(ctx *gin.Context) {
	userID := ctx.Value("userID").(string)
	sessionID := ctx.Value("sessionID").(string)
	err := c.publicService.Logout(userID, sessionID)
	if err != nil {
		logrus.Errorf("Logout: error in logging out err = %v", err)
		responseerror.RespondGenericServerErr(ctx, err, "error in logging out")
		return
	}

	ctx.JSON(http.StatusCreated, "You are logged out")
}
