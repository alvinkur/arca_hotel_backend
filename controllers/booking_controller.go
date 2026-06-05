package controllers

import (
	"net/http"
	"strconv"

	"arca-hotel/config"
	"arca-hotel/models"

	"github.com/gin-gonic/gin"
)

func GetBookings(c *gin.Context) {
	var bookings []models.Booking
	if err := config.DB.Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data booking"})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func CreateBooking(c *gin.Context) {
	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah room ada
	var room models.Room
	if err := config.DB.First(&room, booking.RoomID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room tidak ditemukan"})
		return
	}

	// Cek apakah room tersedia
	if !room.Availability {
		c.JSON(http.StatusConflict, gin.H{"error": "Room sedang tidak tersedia"})
		return
	}

	if err := config.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data booking"})
		return
	}

	// Update status room menjadi tidak tersedia
	config.DB.Model(&room).Update("availability", false)

	c.JSON(http.StatusCreated, booking)
}

func UpdateBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var existing models.Booking
	if err := config.DB.First(&existing, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking tidak ditemukan"})
		return
	}

	var input models.Booking
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.RoomID != 0 {
		var room models.Room
		if err := config.DB.First(&room, input.RoomID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room tidak ditemukan"})
			return
		}
		if !room.Availability && input.RoomID != existing.RoomID {
			c.JSON(http.StatusConflict, gin.H{"error": "Room sedang tidak tersedia"})
			return
		}
	}

	config.DB.Model(&existing).Updates(&input)
	c.JSON(http.StatusOK, existing)
}

func DeleteBooking(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var booking models.Booking
	if err := config.DB.First(&booking, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking tidak ditemukan"})
		return
	}

	config.DB.Model(&models.Room{}).Where("id_room = ?", booking.RoomID).Update("availability", true)
	config.DB.Delete(&booking)
	c.JSON(http.StatusOK, gin.H{"message": "Booking berhasil dihapus"})
}
