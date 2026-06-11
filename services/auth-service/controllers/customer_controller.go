package controllers

import (
	"net/http"
	"strconv"

	"arca-hotel/services/auth-service/config"
	"arca-hotel/services/auth-service/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func GetCustomer(c *gin.Context) {
	var customers []models.Customer
	if err := config.DB.Find(&customers).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data customer"})
		return
	}
	c.JSON(http.StatusOK, customers)
}

func CreateCustomer(c *gin.Context) {
	var customer models.Customer
	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(customer.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal memproses password"})
		return
	}
	customer.Password = string(hashed)

	if err := config.DB.Create(&customer).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data customer"})
		return
	}
	c.JSON(http.StatusCreated, customer)
}

func UpdateCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var existing models.Customer
	if err := config.DB.First(&existing, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer tidak ditemukan"})
		return
	}

	var input models.Customer
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

func DeleteCustomer(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var customer models.Customer
	if err := config.DB.First(&customer, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Customer tidak ditemukan"})
		return
	}

	config.DB.Delete(&customer)
	c.JSON(http.StatusOK, gin.H{"message": "Customer berhasil dihapus"})
}
