package controllers

import (
	"net/http"
	"strconv"

	"arca-hotel/services/room-service/config"
	"arca-hotel/services/room-service/models"

	"github.com/gin-gonic/gin"
)

func GetRoomTypes(c *gin.Context) {
	var types []models.RoomType
	if err := config.DB.Find(&types).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data room type"})
		return
	}
	c.JSON(http.StatusOK, types)
}

func CreateRoomType(c *gin.Context) {
	var rt models.RoomType
	if err := c.ShouldBindJSON(&rt); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&rt).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data room type"})
		return
	}
	c.JSON(http.StatusCreated, rt)
}

func UpdateRoomType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var existing models.RoomType
	if err := config.DB.First(&existing, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room type tidak ditemukan"})
		return
	}

	var input models.RoomType
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&existing).Updates(&input)
	c.JSON(http.StatusOK, existing)
}

func DeleteRoomType(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var rt models.RoomType
	if err := config.DB.First(&rt, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room type tidak ditemukan"})
		return
	}

	config.DB.Delete(&rt)
	c.JSON(http.StatusOK, gin.H{"message": "Room type berhasil dihapus"})
}
