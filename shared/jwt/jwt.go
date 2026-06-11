package jwt

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	jwtlib "github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Name   string `json:"name"`
	jwtlib.RegisteredClaims
}

func getJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "arca-hotel-jwt-secret"
	}
	return []byte(secret)
}

func GenerateToken(userID uint, email, role, name string) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		Email:  email,
		Role:   role,
		Name:   name,
		RegisteredClaims: jwtlib.RegisteredClaims{
			ExpiresAt: jwtlib.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwtlib.NewNumericDate(time.Now()),
		},
	}

	token := jwtlib.NewWithClaims(jwtlib.SigningMethodHS256, claims)
	return token.SignedString(getJWTSecret())
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak ditemukan"})
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwtlib.ParseWithClaims(tokenString, &JWTClaims{}, func(t *jwtlib.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwtlib.SigningMethodHMAC); !ok {
				return nil, jwtlib.ErrSignatureInvalid
			}
			return getJWTSecret(), nil
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			return
		}

		claims, ok := token.Claims.(*JWTClaims)
		if !ok || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token tidak valid"})
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)
		c.Set("name", claims.Name)

		c.Next()
	}
}
