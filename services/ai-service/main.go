package main

import (
	"log"
	"os"

	"arca-hotel/services/ai-service/controllers"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	godotenv.Load("../../.env")
}

func main() {
	r := gin.Default()

	r.POST("/ai-recommend", controllers.RecommendRoom)

	port := os.Getenv("AI_SERVICE_PORT")
	if port == "" {
		port = "8008"
	}
	log.Println("AI Service running on :" + port)
	r.Run(":" + port)
}
