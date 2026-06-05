package controllers

import (
	"net/http"

	"arca-hotel/config"
	"arca-hotel/middleware"
	"arca-hotel/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required,oneof=customer owner staff"`
}

type UserDTO struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type LoginResponse struct {
	Token string  `json:"token"`
	User  UserDTO `json:"user"`
}

func Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var userID uint
	var name, email, hashedPassword string

	switch req.Role {
	case "customer":
		var customer models.Customer
		if err := config.DB.Where("email = ?", req.Email).First(&customer).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
			return
		}
		userID = customer.ID
		name = customer.Name
		email = customer.Email
		hashedPassword = customer.Password

	case "owner":
		var owner models.Owner
		if err := config.DB.Where("email = ?", req.Email).First(&owner).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
			return
		}
		userID = owner.ID
		name = owner.Name
		email = owner.Email
		hashedPassword = owner.Password

	case "staff":
		var staff models.Staff
		if err := config.DB.Where("email = ?", req.Email).First(&staff).Error; err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
			return
		}
		userID = staff.ID
		name = staff.Name
		email = staff.Email
		hashedPassword = staff.Password
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Email atau password salah"})
		return
	}

	token, err := middleware.GenerateToken(userID, email, req.Role, name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal membuat token"})
		return
	}

	c.JSON(http.StatusOK, LoginResponse{
		Token: token,
		User: UserDTO{
			ID:    userID,
			Name:  name,
			Email: email,
			Role:  req.Role,
		},
	})
}
