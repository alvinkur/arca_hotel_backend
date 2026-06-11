package main

import (
	"log"
	"os"

	"arca-hotel/services/booking-service/config"
	"arca-hotel/services/booking-service/controllers"
	"arca-hotel/services/booking-service/models"
	sharedConfig "arca-hotel/shared/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.DB = sharedConfig.ConnectDB("booking",
		&models.Booking{},
	)

	r := gin.Default()

	r.GET("/api/bookings", controllers.GetBookings)
	r.POST("/api/bookings", controllers.CreateBooking)
	r.PUT("/api/bookings/:id", controllers.UpdateBooking)
	r.DELETE("/api/bookings/:id", controllers.DeleteBooking)

	port := os.Getenv("BOOKING_SERVICE_PORT")
	if port == "" {
		port = "8003"
	}
	log.Println("Booking Service running on :" + port)
	r.Run(":" + port)
}
