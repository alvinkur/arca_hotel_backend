package main

import (
	"log"
	"os"

	"arca-hotel/services/room-service/config"
	"arca-hotel/services/room-service/controllers"
	"arca-hotel/services/room-service/models"
	sharedConfig "arca-hotel/shared/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.DB = sharedConfig.ConnectDB("room",
		&models.RoomType{},
		&models.Room{},
	)

	r := gin.Default()

	r.GET("/api/room-types", controllers.GetRoomTypes)
	r.POST("/api/room-types", controllers.CreateRoomType)
	r.PUT("/api/room-types/:id", controllers.UpdateRoomType)
	r.DELETE("/api/room-types/:id", controllers.DeleteRoomType)

	r.GET("/api/rooms", controllers.GetRooms)
	r.POST("/api/rooms", controllers.CreateRoom)
	r.PUT("/api/rooms/:id", controllers.UpdateRoom)
	r.DELETE("/api/rooms/:id", controllers.DeleteRoom)

	port := os.Getenv("ROOM_SERVICE_PORT")
	if port == "" {
		port = "8002"
	}
	log.Println("Room Service running on :" + port)
	r.Run(":" + port)
}
