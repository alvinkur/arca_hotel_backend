package main

import (
	"log"
	"os"

	"arca-hotel/services/auth-service/config"
	"arca-hotel/services/auth-service/controllers"
	"arca-hotel/services/auth-service/models"
	sharedConfig "arca-hotel/shared/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.DB = sharedConfig.ConnectDB("auth",
		&models.Customer{},
		&models.Owner{},
		&models.Staff{},
	)

	r := gin.Default()

	r.POST("/api/login", controllers.Login)
	r.GET("/api/customers", controllers.GetCustomer)
	r.POST("/api/customers", controllers.CreateCustomer)
	r.PUT("/api/customers/:id", controllers.UpdateCustomer)
	r.DELETE("/api/customers/:id", controllers.DeleteCustomer)

	r.GET("/api/owners", controllers.GetOwner)
	r.POST("/api/owners", controllers.CreateOwner)
	r.PUT("/api/owners/:id", controllers.UpdateOwner)
	r.DELETE("/api/owners/:id", controllers.DeleteOwner)

	r.GET("/api/staffs", controllers.GetStaffs)
	r.POST("/api/staffs", controllers.CreateStaff)
	r.PUT("/api/staffs/:id", controllers.UpdateStaff)
	r.DELETE("/api/staffs/:id", controllers.DeleteStaff)

	port := os.Getenv("AUTH_SERVICE_PORT")
	if port == "" {
		port = "8001"
	}
	log.Println("Auth Service running on :" + port)
	r.Run(":" + port)
}
