package main

import (
	"log"
	"os"

	"arca-hotel/services/payment-service/config"
	"arca-hotel/services/payment-service/controllers"
	"arca-hotel/services/payment-service/models"
	sharedConfig "arca-hotel/shared/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.DB = sharedConfig.ConnectDB("payment",
		&models.Payment{},
	)

	r := gin.Default()

	r.GET("/api/payments", controllers.GetPayments)
	r.POST("/api/payments", controllers.CreatePayment)
	r.PUT("/api/payments/:id", controllers.UpdatePayment)
	r.DELETE("/api/payments/:id", controllers.DeletePayment)

	port := os.Getenv("PAYMENT_SERVICE_PORT")
	if port == "" {
		port = "8004"
	}
	log.Println("Payment Service running on :" + port)
	r.Run(":" + port)
}
