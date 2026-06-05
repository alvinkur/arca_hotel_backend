package controllers

import (
	"net/http"
	"strconv"

	"arca-hotel/config"
	"arca-hotel/models"

	"github.com/gin-gonic/gin"
)

func GetChats(c *gin.Context) {
	var chats []models.Chat
	if err := config.DB.Preload("ChatMessages").Find(&chats).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data chat"})
		return
	}
	c.JSON(http.StatusOK, chats)
}

func CreateChat(c *gin.Context) {
	var chat models.Chat
	if err := c.ShouldBindJSON(&chat); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&chat).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data chat"})
		return
	}
	c.JSON(http.StatusCreated, chat)
}

func UpdateChat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var existing models.Chat
	if err := config.DB.First(&existing, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat tidak ditemukan"})
		return
	}

	var input models.Chat
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&existing).Updates(&input)
	c.JSON(http.StatusOK, existing)
}

func DeleteChat(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var chat models.Chat
	if err := config.DB.First(&chat, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Chat tidak ditemukan"})
		return
	}

	config.DB.Where("id_chat = ?", id).Delete(&models.ChatMessage{})
	config.DB.Delete(&chat)
	c.JSON(http.StatusOK, gin.H{"message": "Chat berhasil dihapus"})
}
