package controllers

import (
	"net/http"
	"os"
	"strconv"

	"arca-hotel/services/payment-service/clients"
	"arca-hotel/services/payment-service/config"
	"arca-hotel/services/payment-service/models"

	"github.com/gin-gonic/gin"
)

var bookingClient *clients.BookingClient

func initClient() {
	if bookingClient == nil {
		base := os.Getenv("BOOKING_SERVICE_URL")
		if base == "" {
			base = "http://localhost:8003"
		}
		bookingClient = clients.NewBookingClient(base)
	}
}

func GetPayments(c *gin.Context) {
	var payments []models.Payment
	if err := config.DB.Find(&payments).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data payment"})
		return
	}
	c.JSON(http.StatusOK, payments)
}

func CreatePayment(c *gin.Context) {
	initClient()

	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := bookingClient.ValidateBooking(payment.BookingID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Booking tidak ditemukan"})
		return
	}

	if err := config.DB.Create(&payment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data payment"})
		return
	}

	bookingClient.UpdatePaymentStatus(payment.BookingID, "paid")

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
		initClient()
		if err := bookingClient.ValidateBooking(input.BookingID); err != nil {
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
