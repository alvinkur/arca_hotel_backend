package main

import (
	"log"
	"os"

	"arca-hotel/services/review-service/config"
	"arca-hotel/services/review-service/controllers"
	"arca-hotel/services/review-service/models"
	sharedConfig "arca-hotel/shared/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.DB = sharedConfig.ConnectDB("review",
		&models.Review{},
	)

	r := gin.Default()

	r.GET("/api/reviews", controllers.GetReviews)
	r.POST("/api/reviews", controllers.CreateReview)
	r.PUT("/api/reviews/:id", controllers.UpdateReview)
	r.DELETE("/api/reviews/:id", controllers.DeleteReview)

	port := os.Getenv("REVIEW_SERVICE_PORT")
	if port == "" {
		port = "8006"
	}
	log.Println("Review Service running on :" + port)
	r.Run(":" + port)
}
