package controllers

import (
	"net/http"
	"strconv"

	"arca-hotel/config"
	"arca-hotel/models"

	"github.com/gin-gonic/gin"
)

func GetRevenueReports(c *gin.Context) {
	var reports []models.RevenueReport
	if err := config.DB.Find(&reports).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data revenue report"})
		return
	}
	c.JSON(http.StatusOK, reports)
}

func CreateRevenueReport(c *gin.Context) {
	var report models.RevenueReport
	if err := c.ShouldBindJSON(&report); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := config.DB.Create(&report).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data revenue report"})
		return
	}
	c.JSON(http.StatusCreated, report)
}

func UpdateRevenueReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var existing models.RevenueReport
	if err := config.DB.First(&existing, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Revenue report tidak ditemukan"})
		return
	}

	var input models.RevenueReport
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	config.DB.Model(&existing).Updates(&input)
	c.JSON(http.StatusOK, existing)
}

func DeleteRevenueReport(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var report models.RevenueReport
	if err := config.DB.First(&report, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Revenue report tidak ditemukan"})
		return
	}

	config.DB.Delete(&report)
	c.JSON(http.StatusOK, gin.H{"message": "Revenue report berhasil dihapus"})
}
