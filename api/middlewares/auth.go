package middlewares

import (
	"errors"
	"github.com/Shresth92/audiophile/internal"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/responseerror"
	"github.com/Shresth92/audiophile/services"
	"github.com/Shresth92/audiophile/utils"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"time"
)

type AuthMiddleware struct {
	handler     *internal.RequestHandler
	authService services.PublicService
}

func NewAuthMiddleware(
	handler *internal.RequestHandler,
	authService services.PublicService,
) *AuthMiddleware {
	return &AuthMiddleware{
		handler:     handler,
		authService: authService,
	}
}

func (m *AuthMiddleware) Setup(ctx *gin.Context) {
	token := ctx.Request.Header.Get("authorization")
	claims := &models.Claims{}
	if token == "" {
		responseerror.RespondClientErr(ctx, errors.New("token not sent in header"), http.StatusUnauthorized, "token not sent in header")
		ctx.Abort()
		return
	} else {
		parseToken, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(utils.GetEnvValue("JwtKey")), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				responseerror.RespondClientErr(ctx, err, http.StatusUnauthorized, "error in validating token")
				ctx.Abort()
				return
			}
			responseerror.RespondClientErr(ctx, err, http.StatusUnauthorized, "expired token")
			ctx.Abort()
			return
		}
		if !parseToken.Valid {
			responseerror.RespondClientErr(ctx, errors.New("token is not valid"), http.StatusUnauthorized, "token is not valid")
			ctx.Abort()
			return
		}

		sessionId := claims.SessionId
		userId := claims.UserId

		sessionEndTime, err := m.authService.CheckSession(sessionId, userId)
		if err != nil {
			responseerror.RespondClientErr(ctx, err, http.StatusInternalServerError, "error in checking session")
			ctx.Abort()
			return
		}

		currentTime := time.Now()
		if sessionEndTime.Before(currentTime) {
			responseerror.RespondClientErr(ctx, err, http.StatusUnauthorized, "already logged out")
			ctx.Abort()
			return
		}
		ctx.Set("userID", claims.UserId)
		ctx.Set("sessionID", claims.SessionId)
		ctx.Set("role", claims.Role)
		ctx.Next()
	}
}
