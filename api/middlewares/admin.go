package middlewares

import (
	"errors"
	"github.com/Shresth92/audiophile/internal"
	"github.com/Shresth92/audiophile/models"
	"github.com/Shresth92/audiophile/responseerror"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AdminMiddleware struct {
	handler *internal.RequestHandler
}

func NewAdminMiddleware(
	handler *internal.RequestHandler,
) *AdminMiddleware {
	return &AdminMiddleware{
		handler: handler,
	}
}

func (m *AdminMiddleware) Setup(ctx *gin.Context) {
	role := ctx.Value("role")
	if role == models.Admin {
		ctx.Next()
	} else {
		responseerror.RespondClientErr(ctx, errors.New("not admin"), http.StatusUnauthorized, "not admin")
		ctx.Abort()
		return
	}
}
