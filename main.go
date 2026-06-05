package main

import (
	"arca-hotel/config"
	"arca-hotel/controllers"
	"arca-hotel/middleware"

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

	// Public Routes
	r.POST("/api/login", controllers.Login)

	// Protected API Routes
	api := r.Group("/api")
	api.Use(middleware.AuthMiddleware())
	{
		api.GET("/customers", controllers.GetCustomer)
		api.POST("/customers", controllers.CreateCustomer)
		api.PUT("/customers/:id", controllers.UpdateCustomer)
		api.DELETE("/customers/:id", controllers.DeleteCustomer)

		api.GET("/owners", controllers.GetOwner)
		api.POST("/owners", controllers.CreateOwner)
		api.PUT("/owners/:id", controllers.UpdateOwner)
		api.DELETE("/owners/:id", controllers.DeleteOwner)

		api.GET("/room-types", controllers.GetRoomTypes)
		api.POST("/room-types", controllers.CreateRoomType)
		api.PUT("/room-types/:id", controllers.UpdateRoomType)
		api.DELETE("/room-types/:id", controllers.DeleteRoomType)

		api.GET("/rooms", controllers.GetRooms)
		api.POST("/rooms", controllers.CreateRoom)
		api.PUT("/rooms/:id", controllers.UpdateRoom)
		api.DELETE("/rooms/:id", controllers.DeleteRoom)

		api.GET("/bookings", controllers.GetBookings)
		api.POST("/bookings", controllers.CreateBooking)
		api.PUT("/bookings/:id", controllers.UpdateBooking)
		api.DELETE("/bookings/:id", controllers.DeleteBooking)

		api.GET("/staffs", controllers.GetStaffs)
		api.POST("/staffs", controllers.CreateStaff)
		api.PUT("/staffs/:id", controllers.UpdateStaff)
		api.DELETE("/staffs/:id", controllers.DeleteStaff)

		api.GET("/payments", controllers.GetPayments)
		api.POST("/payments", controllers.CreatePayment)
		api.PUT("/payments/:id", controllers.UpdatePayment)
		api.DELETE("/payments/:id", controllers.DeletePayment)

		api.GET("/reviews", controllers.GetReviews)
		api.POST("/reviews", controllers.CreateReview)
		api.PUT("/reviews/:id", controllers.UpdateReview)
		api.DELETE("/reviews/:id", controllers.DeleteReview)

		api.GET("/chats", controllers.GetChats)
		api.POST("/chats", controllers.CreateChat)
		api.PUT("/chats/:id", controllers.UpdateChat)
		api.DELETE("/chats/:id", controllers.DeleteChat)
		api.GET("/chats/:id/messages", controllers.GetChatMessages)
		api.POST("/chats/:id/messages", controllers.CreateChatMessage)
		api.PUT("/chats/:id/messages/:msgId", controllers.UpdateChatMessage)
		api.DELETE("/chats/:id/messages/:msgId", controllers.DeleteChatMessage)

		api.GET("/revenue_reports", controllers.GetRevenueReports)
		api.POST("/revenue_reports", controllers.CreateRevenueReport)
		api.PUT("/revenue_reports/:id", controllers.UpdateRevenueReport)
		api.DELETE("/revenue_reports/:id", controllers.DeleteRevenueReport)
	}

	r.Run(":8080")
}
