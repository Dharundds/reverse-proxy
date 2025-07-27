package handlers

import (
	"net/http"
	"reverse-proxy/internal/constants"
	"reverse-proxy/internal/helpers"

	"github.com/gin-gonic/gin"
)

func RPHandler(c *gin.Context) {
	host := c.Request.Host

	val, ok := constants.RPCtxManager.GetContext()[host]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "Reverse proxy not found"})
		return
	}

	targetHost := "http://" + host + ":" + val

	proxy := helpers.CreateReverseProxy(targetHost)
	proxy(c)
}

func AddRPHandler(c *gin.Context) {
	type body struct {
		DN   string `json:"domainName"`
		Port string `json:"port"`
	}

	var payload body
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	ok := helpers.AddRP(payload.DN, payload.Port)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to add reverse proxy"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Reverse proxy added successfully"})
}

func DelRPHandler(c *gin.Context) {
	type body struct {
		DN string `json:"domainName"`
	}

	var payload body
	if err := c.BindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}

	ok := helpers.RemoveRP(payload.DN)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to delete reverse proxy"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Reverse proxy deleted successfully"})
}

func GetRPHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Reverse proxy fetched successfully", "data": constants.RPCtxManager.GetContext()})
}

func ReloadRPHandler(c *gin.Context) {
	res := helpers.LoadRedisContext()
	if res == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Failed to reload reverse proxy"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reverse proxy reloaded successfully", "data": res})
}
