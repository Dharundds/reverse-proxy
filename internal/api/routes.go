package api

import (
	"reverse-proxy/internal/api/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterAPIRoutes(g *gin.Engine) {
	// g.GET("/*path", handlers.RPHandler)

	apiGroup := g.Group("/api")
	apiGroup.POST("/rp", handlers.AddRPHandler)
	apiGroup.DELETE("/rp", handlers.DelRPHandler)
	apiGroup.GET("/rp", handlers.GetRPHandler)
	apiGroup.GET("/rp/reload", handlers.ReloadRPHandler)
}

func RegisterRPHandler(g *gin.Engine) {
	g.Any("/*path", handlers.RPHandler)
}
