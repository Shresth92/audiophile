package middlewares

import (
	"errors"
	"github.com/Shresth92/audiophile/internal"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/responseerror"
	"github.com/gin-gonic/gin"
	"net/http"
)

type UserMiddleware struct {
	handler *internal.RequestHandler
}

func NewUserMiddleware(
	handler *internal.RequestHandler,
) *UserMiddleware {
	return &UserMiddleware{
		handler: handler,
	}
}

func (m *UserMiddleware) Setup(ctx *gin.Context) {
	role := ctx.Value("role")
	if role == models.User {
		ctx.Next()
	} else {
		responseerror.RespondClientErr(ctx, errors.New("not user"), http.StatusUnauthorized, "not user")
		ctx.Abort()
		return
	}
}
