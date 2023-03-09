package internal

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type RequestHandler struct {
	Gin *gin.Engine
}

// NewRequestHandler creates a new request handler
func NewRequestHandler() *RequestHandler {
	engine := gin.New()
	engine.ForwardedByClientIP = true
	engine.Use(gin.Recovery())
	engine.Use(func(c *gin.Context) {
		c.Header("Content-Type", "application/json")
		c.Next()
	})
	engine.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"responseerror": "Page not found",
		})
	})
	return &RequestHandler{Gin: engine}
}
