package controllers

import (
	"net/http"
	"strconv"

	"arca-hotel/config"
	"arca-hotel/models"

	"github.com/gin-gonic/gin"
)

func GetReviews(c *gin.Context) {
	var reviews []models.Review
	if err := config.DB.Find(&reviews).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data review"})
		return
	}
	c.JSON(http.StatusOK, reviews)
}

func CreateReview(c *gin.Context) {
	var review models.Review
	if err := c.ShouldBindJSON(&review); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validasi rating 1-5
	if review.Rating < 1 || review.Rating > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating harus antara 1 sampai 5"})
		return
	}

	if err := config.DB.Create(&review).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal menyimpan data review"})
		return
	}
	c.JSON(http.StatusCreated, review)
}

func UpdateReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var existing models.Review
	if err := config.DB.First(&existing, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review tidak ditemukan"})
		return
	}

	var input models.Review
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Rating != 0 && (input.Rating < 1 || input.Rating > 5) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Rating harus antara 1 sampai 5"})
		return
	}

	config.DB.Model(&existing).Updates(&input)
	c.JSON(http.StatusOK, existing)
}

func DeleteReview(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID tidak valid"})
		return
	}

	var review models.Review
	if err := config.DB.First(&review, uint(id)).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Review tidak ditemukan"})
		return
	}

	config.DB.Delete(&review)
	c.JSON(http.StatusOK, gin.H{"message": "Review berhasil dihapus"})
}
