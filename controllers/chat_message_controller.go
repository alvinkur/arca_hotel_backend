package controllers

import (
	"net/http"
	"strconv"

	"arca-hotel/config"
	"arca-hotel/models"

	"github.com/gin-gonic/gin"
)

func GetChatMessages(c *gin.Context) {
	chatID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID chat tidak valid"})
		return
	}

	var messages []models.ChatMessage
	if err := config.DB.Where("id_chat = ?", chatID).Order("date ASC").Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data pesan"})
		return
	}
	c.JSON(http.StatusOK, messages)
}

func CreateChatMessage(c *gin.Context) {
	chatID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID chat tidak valid"})
		return
	}

	// Cek chat ada
	var chat models.Chat
	if err := config.DB.First(&chat, chatID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat tidak ditemukan"})
		return
	}

	var msg models.ChatMessage
	if err := c.ShouldBindJSON(&msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg.ChatID = uint(chatID)

	if err := config.DB.Create(&msg).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan pesan"})
		return
	}
	c.JSON(http.StatusCreated, msg)
}

func UpdateChatMessage(c *gin.Context) {
	msgID, err := strconv.Atoi(c.Param("msgId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pesan tidak valid"})
		return
	}

	var existing models.ChatMessage
	if err := config.DB.First(&existing, uint(msgID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesan tidak ditemukan"})
		return
	}

	var input models.ChatMessage
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&existing).Updates(&input)
	c.JSON(http.StatusOK, existing)
}

func DeleteChatMessage(c *gin.Context) {
	msgID, err := strconv.Atoi(c.Param("msgId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID pesan tidak valid"})
		return
	}

	var msg models.ChatMessage
	if err := config.DB.First(&msg, uint(msgID)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Pesan tidak ditemukan"})
		return
	}

	config.DB.Delete(&msg)
	c.JSON(http.StatusOK, gin.H{"message": "Pesan berhasil dihapus"})
}
