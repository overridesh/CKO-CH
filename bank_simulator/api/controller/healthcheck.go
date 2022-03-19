package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func HealthcheckController(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"ok": true,
	})
}
