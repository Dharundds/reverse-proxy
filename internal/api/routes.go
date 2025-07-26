package api

import "github.com/gin-gonic/gin"

func HandleFuncs(g *gin.Engine) {
	g.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello, World!"})
	})
}
