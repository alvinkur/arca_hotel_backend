package main

import (
	"log"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"arca-hotel/shared/jwt"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	authURL := getEnv("AUTH_SERVICE_URL", "http://localhost:8001")
	roomURL := getEnv("ROOM_SERVICE_URL", "http://localhost:8002")
	bookingURL := getEnv("BOOKING_SERVICE_URL", "http://localhost:8003")
	paymentURL := getEnv("PAYMENT_SERVICE_URL", "http://localhost:8004")
	chatURL := getEnv("CHAT_SERVICE_URL", "http://localhost:8005")
	reviewURL := getEnv("REVIEW_SERVICE_URL", "http://localhost:8006")
	reportURL := getEnv("REPORT_SERVICE_URL", "http://localhost:8007")
	aiURL := getEnv("AI_SERVICE_URL", "http://localhost:8008")

	r := gin.Default()

	// CORS
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

	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Public pages (no auth needed for the page itself)
	r.GET("/login", func(c *gin.Context) { c.HTML(200, "login.html", nil) })

	// Dashboard & entity pages
	r.GET("/", func(c *gin.Context) { c.HTML(200, "dashboard.html", nil) })
	r.GET("/dashboard", func(c *gin.Context) { c.HTML(200, "dashboard.html", nil) })
	r.GET("/customer", func(c *gin.Context) { c.HTML(200, "customer.html", nil) })
	r.GET("/owner", func(c *gin.Context) { c.HTML(200, "owner.html", nil) })
	r.GET("/staff", func(c *gin.Context) { c.HTML(200, "staff.html", nil) })
	r.GET("/room", func(c *gin.Context) { c.HTML(200, "room.html", nil) })
	r.GET("/room-type", func(c *gin.Context) { c.HTML(200, "room_type.html", nil) })
	r.GET("/booking", func(c *gin.Context) { c.HTML(200, "booking.html", nil) })
	r.GET("/payment", func(c *gin.Context) { c.HTML(200, "payment.html", nil) })
	r.GET("/chat", func(c *gin.Context) { c.HTML(200, "chat.html", nil) })
	r.GET("/review", func(c *gin.Context) { c.HTML(200, "review.html", nil) })
	r.GET("/revenue-report", func(c *gin.Context) { c.HTML(200, "revenue_report.html", nil) })

	// Public routes (no auth)
	r.POST("/api/login", proxy(authURL))
	r.POST("/ai-recommend", proxy(aiURL))

	// Protected routes
	api := r.Group("/api")
	api.Use(jwt.AuthMiddleware())
	{
		api.Any("/customers/*path", proxy(authURL))
		api.Any("/owners/*path", proxy(authURL))
		api.Any("/staffs/*path", proxy(authURL))

		api.Any("/room-types/*path", proxy(roomURL))
		api.Any("/rooms/*path", proxy(roomURL))

		api.Any("/bookings/*path", proxy(bookingURL))

		api.Any("/payments/*path", proxy(paymentURL))

		api.Any("/chats/*path", proxy(chatURL))

		api.Any("/reviews/*path", proxy(reviewURL))

		api.Any("/revenue_reports/*path", proxy(reportURL))
	}

	// Also handle exact-path matches (no /*path suffix)
	api.GET("/customers", proxy(authURL))
	api.POST("/customers", proxy(authURL))
	api.GET("/owners", proxy(authURL))
	api.POST("/owners", proxy(authURL))
	api.GET("/staffs", proxy(authURL))
	api.POST("/staffs", proxy(authURL))
	api.GET("/room-types", proxy(roomURL))
	api.POST("/room-types", proxy(roomURL))
	api.GET("/rooms", proxy(roomURL))
	api.POST("/rooms", proxy(roomURL))
	api.GET("/bookings", proxy(bookingURL))
	api.POST("/bookings", proxy(bookingURL))
	api.GET("/payments", proxy(paymentURL))
	api.POST("/payments", proxy(paymentURL))
	api.GET("/chats", proxy(chatURL))
	api.POST("/chats", proxy(chatURL))
	api.GET("/reviews", proxy(reviewURL))
	api.POST("/reviews", proxy(reviewURL))
	api.GET("/revenue_reports", proxy(reportURL))
	api.POST("/revenue_reports", proxy(reportURL))

	port := os.Getenv("GATEWAY_PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("API Gateway running on :" + port)
	r.Run(":" + port)
}

func proxy(target string) gin.HandlerFunc {
	targetURL, _ := url.Parse(target)
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	return func(c *gin.Context) {
		// Strip the prefix that was used for routing so the backend gets a clean path.
		// The reverse proxy will rewrite the path to match the target URL.
		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func init() {
	godotenv.Load(".env")

	// Ensure we don't have trailing slashes on base URLs
	cleanEnv := func(key string) {
		if v := os.Getenv(key); v != "" {
			os.Setenv(key, strings.TrimRight(v, "/"))
		}
	}
	for _, k := range []string{
		"AUTH_SERVICE_URL", "ROOM_SERVICE_URL", "BOOKING_SERVICE_URL",
		"PAYMENT_SERVICE_URL", "CHAT_SERVICE_URL", "REVIEW_SERVICE_URL",
		"REPORT_SERVICE_URL", "AI_SERVICE_URL",
	} {
		cleanEnv(k)
	}
}
