package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-org/oj-platform/pkg/response"
)

func HealthCheck(c *gin.Context) {
	response.Success(c, gin.H{
		"status": "ok",
	})
}

func Welcome(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to OJ Platform API",
		"version": "1.0.0",
	})
}
