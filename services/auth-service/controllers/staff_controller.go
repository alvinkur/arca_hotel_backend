package controllers

import (
	"net/http"
	"strconv"

	"arca-hotel/services/auth-service/config"
	"arca-hotel/services/auth-service/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetStaffs(c *gin.Context) {
	var staffs []models.Staff
	if err := config.DB.Find(&staffs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data staff"})
		return
	}
	c.JSON(http.StatusOK, staffs)
}

func CreateStaff(c *gin.Context) {
	var staff models.Staff
	if err := c.ShouldBindJSON(&staff); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(staff.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}
	staff.Password = string(hashed)

	if err := config.DB.Create(&staff).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data staff"})
		return
	}
	c.JSON(http.StatusCreated, staff)
}

func UpdateStaff(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var existing models.Staff
	if err := config.DB.First(&existing, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff tidak ditemukan"})
		return
	}

	var input models.Staff
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Password != "" {
		hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		input.Password = string(hashed)
	}

	config.DB.Model(&existing).Updates(&input)
	c.JSON(http.StatusOK, existing)
}

func DeleteStaff(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var staff models.Staff
	if err := config.DB.First(&staff, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Staff tidak ditemukan"})
		return
	}

	config.DB.Delete(&staff)
	c.JSON(http.StatusOK, gin.H{"message": "Staff berhasil dihapus"})
}
