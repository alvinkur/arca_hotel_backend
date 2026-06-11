package main

import (
	"log"
	"os"

	"arca-hotel/services/report-service/config"
	"arca-hotel/services/report-service/controllers"
	"arca-hotel/services/report-service/models"
	sharedConfig "arca-hotel/shared/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.DB = sharedConfig.ConnectDB("report",
		&models.RevenueReport{},
	)

	r := gin.Default()

	r.GET("/api/revenue_reports", controllers.GetRevenueReports)
	r.POST("/api/revenue_reports", controllers.CreateRevenueReport)
	r.PUT("/api/revenue_reports/:id", controllers.UpdateRevenueReport)
	r.DELETE("/api/revenue_reports/:id", controllers.DeleteRevenueReport)

	port := os.Getenv("REPORT_SERVICE_PORT")
	if port == "" {
		port = "8007"
	}
	log.Println("Report Service running on :" + port)
	r.Run(":" + port)
}
