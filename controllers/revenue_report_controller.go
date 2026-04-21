package controllers

import (
	"net/http"

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
