package controllers

import (
	"net/http"
	"strconv"

	"arca-hotel/config"
	"arca-hotel/models"

	"github.com/gin-gonic/gin"
)

func GetPayments(c *gin.Context) {
	var payments []models.Payment
	if err := config.DB.Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data payment"})
		return
	}
	c.JSON(http.StatusOK, payments)
}

func CreatePayment(c *gin.Context) {
	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Cek apakah booking ada
	var booking models.Booking
	if err := config.DB.First(&booking, payment.BookingID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking tidak ditemukan"})
		return
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data payment"})
		return
	}

	// Update status payment di booking
	config.DB.Model(&booking).Update("status_payment", "paid")

	c.JSON(http.StatusCreated, payment)
}

func UpdatePayment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var existing models.Payment
	if err := config.DB.First(&existing, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment tidak ditemukan"})
		return
	}

	var input models.Payment
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.BookingID != 0 {
		var booking models.Booking
		if err := config.DB.First(&booking, input.BookingID).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Booking tidak ditemukan"})
			return
		}
	}

	config.DB.Model(&existing).Updates(&input)
	c.JSON(http.StatusOK, existing)
}

func DeletePayment(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var payment models.Payment
	if err := config.DB.First(&payment, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Payment tidak ditemukan"})
		return
	}

	config.DB.Delete(&payment)
	c.JSON(http.StatusOK, gin.H{"message": "Payment berhasil dihapus"})
}
