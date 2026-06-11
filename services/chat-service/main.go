package main

import (
	"log"
	"os"

	"arca-hotel/services/chat-service/config"
	"arca-hotel/services/chat-service/controllers"
	"arca-hotel/services/chat-service/models"
	sharedConfig "arca-hotel/shared/config"

	"github.com/gin-gonic/gin"
)

func main() {
	config.DB = sharedConfig.ConnectDB("chat",
		&models.Chat{},
		&models.ChatMessage{},
	)

	r := gin.Default()

	r.GET("/api/chats", controllers.GetChats)
	r.POST("/api/chats", controllers.CreateChat)
	r.PUT("/api/chats/:id", controllers.UpdateChat)
	r.DELETE("/api/chats/:id", controllers.DeleteChat)
	r.GET("/api/chats/:id/messages", controllers.GetChatMessages)
	r.POST("/api/chats/:id/messages", controllers.CreateChatMessage)
	r.PUT("/api/chats/:id/messages/:msgId", controllers.UpdateChatMessage)
	r.DELETE("/api/chats/:id/messages/:msgId", controllers.DeleteChatMessage)

	port := os.Getenv("CHAT_SERVICE_PORT")
	if port == "" {
		port = "8005"
	}
	log.Println("Chat Service running on :" + port)
	r.Run(":" + port)
}
