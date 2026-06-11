package controllers

import (
	"net/http"
	"strconv"

	"arca-hotel/services/room-service/config"
	"arca-hotel/services/room-service/models"

	"github.com/gin-gonic/gin"
)

func GetRooms(c *gin.Context) {
	var rooms []models.Room
	if err := config.DB.Preload("RoomType").Find(&rooms).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data room"})
		return
	}
	c.JSON(http.StatusOK, rooms)
}

func CreateRoom(c *gin.Context) {
	var room models.Room
	if err := c.ShouldBindJSON(&room); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var rt models.RoomType
	if err := config.DB.First(&rt, room.RoomTypeID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room type tidak ditemukan"})
		return
	}

	if err := config.DB.Create(&room).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data room"})
		return
	}

	config.DB.Preload("RoomType").First(&room, room.ID)
	c.JSON(http.StatusCreated, room)
}

func UpdateRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var existing models.Room
	if err := config.DB.First(&existing, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room tidak ditemukan"})
		return
	}

	var input models.Room
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.RoomTypeID != 0 {
		var rt models.RoomType
		if err := config.DB.First(&rt, input.RoomTypeID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room type tidak ditemukan"})
			return
		}
	}

	config.DB.Model(&existing).Updates(&input)
	config.DB.Preload("RoomType").First(&existing, uint(id))
	c.JSON(http.StatusOK, existing)
}

func DeleteRoom(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var room models.Room
	if err := config.DB.First(&room, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room tidak ditemukan"})
		return
	}

	config.DB.Delete(&room)
	c.JSON(http.StatusOK, gin.H{"message": "Room berhasil dihapus"})
}
