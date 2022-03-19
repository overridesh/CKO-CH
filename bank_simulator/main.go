package main

import (
	"github.com/gin-gonic/gin"

	controller "github.com/overridesh/checkout-challenger/bank_simulator/api/controller"
)

// PORT static port
const PORT string = ":80"

func main() {
	gin.SetMode(gin.ReleaseMode)

	// Get router
	router := NewRouter()

	// Run server
	router.Run(PORT)
}

// NewRouter create a new router with routes
func NewRouter() *gin.Engine {
	// Router with gin
	router := gin.Default()

	router.GET("/healthcheck", controller.HealthcheckController)

	router.POST("/transactions", controller.TransactionControllerFunc)

	return router
}
