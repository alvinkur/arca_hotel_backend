package controllers

import (
	"net/http"
	"os"
	"strconv"

	"arca-hotel/services/booking-service/clients"
	"arca-hotel/services/booking-service/config"
	"arca-hotel/services/booking-service/models"

	"github.com/gin-gonic/gin"
)

var roomClient *clients.RoomClient
var authClient *clients.AuthClient

func initClients() {
	if roomClient == nil {
		base := os.Getenv("ROOM_SERVICE_URL")
		if base == "" {
			base = "http://localhost:8002"
		}
		roomClient = clients.NewRoomClient(base)
	}
	if authClient == nil {
		base := os.Getenv("AUTH_SERVICE_URL")
		if base == "" {
			base = "http://localhost:8001"
		}
		authClient = clients.NewAuthClient(base)
	}
}

func GetBookings(c *gin.Context) {
	var bookings []models.Booking
	if err := config.DB.Find(&bookings).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data booking"})
		return
	}
	c.JSON(http.StatusOK, bookings)
}

func CreateBooking(c *gin.Context) {
	initClients()

	var booking models.Booking
	if err := c.ShouldBindJSON(&booking); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := authClient.ValidateCustomer(booking.CustomerID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer tidak ditemukan"})
		return
	}

	room, err := roomClient.GetRoom(booking.RoomID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Room tidak ditemukan"})
		return
	}
	if !room.Availability {
		c.JSON(http.StatusConflict, gin.H{"error": "Room sedang tidak tersedia"})
		return
	}

	if err := config.DB.Create(&booking).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data booking"})
		return
	}

	roomClient.SetAvailability(booking.RoomID, false)

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

	if input.RoomID != 0 && input.RoomID != existing.RoomID {
		initClients()
		room, err := roomClient.GetRoom(input.RoomID)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Room tidak ditemukan"})
			return
		}
		if !room.Availability {
			c.JSON(http.StatusConflict, gin.H{"error": "Room sedang tidak tersedia"})
			return
		}
		// Release old room, claim new room
		roomClient.SetAvailability(existing.RoomID, true)
		roomClient.SetAvailability(input.RoomID, false)
	}

	config.DB.Model(&existing).Updates(&input)
	c.JSON(http.StatusOK, existing)
}

func DeleteBooking(c *gin.Context) {
	initClients()

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

	// Release room
	roomClient.SetAvailability(booking.RoomID, true)
	config.DB.Delete(&booking)
	c.JSON(http.StatusOK, gin.H{"message": "Booking berhasil dihapus"})
}
