package controllers

import (
	"net/http"

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
