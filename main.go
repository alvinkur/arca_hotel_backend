package main

import (
	"arca-hotel/config"
	"arca-hotel/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.ConnectDB()

	r := gin.Default()

	// CORS Middleware — izinkan frontend mengakses API
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	r.LoadHTMLGlob("templates/*")

	// HTML Pages
	r.GET("/", func(c *gin.Context) { c.HTML(200, "dashboard.html", nil) })
	r.GET("/customer", func(c *gin.Context) { c.HTML(200, "customer.html", nil) })
	r.GET("/owner", func(c *gin.Context) { c.HTML(200, "owner.html", nil) })
	r.GET("/room", func(c *gin.Context) { c.HTML(200, "room.html", nil) })

	// API Routes
	api := r.Group("/api")
	{
		api.GET("/customers", controllers.GetCustomer)
		api.POST("/customers", controllers.CreateCustomer)

		api.GET("/owners", controllers.GetOwner)
		api.POST("/owners", controllers.CreateOwner)

		api.GET("/rooms", controllers.GetRooms)
		api.POST("/rooms", controllers.CreateRoom)

		api.GET("/bookings", controllers.GetBookings)
		api.POST("/bookings", controllers.CreateBooking)

		api.GET("/staffs", controllers.GetStaffs)
		api.POST("/staffs", controllers.CreateStaff)

		api.GET("/payments", controllers.GetPayments)
		api.POST("/payments", controllers.CreatePayment)

		api.GET("/reviews", controllers.GetReviews)
		api.POST("/reviews", controllers.CreateReview)

		api.GET("/chats", controllers.GetChats)
		api.POST("/chats", controllers.CreateChat)

		api.GET("/revenue_reports", controllers.GetRevenueReports)
		api.POST("/revenue_reports", controllers.CreateRevenueReport)
	}

	r.Run(":8080")
}
