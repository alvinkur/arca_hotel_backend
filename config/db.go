package config

import (
	"fmt"
	"log"
	"os"

	"arca-hotel/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	host     := getEnv("DB_HOST", "localhost")
	user     := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "")
	dbname   := getEnv("DB_NAME", "db_hotel_arca")
	port     := getEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal connect database: " + err.Error())
	}

	errMigrate := db.AutoMigrate(
		&models.RoomType{},
		&models.Room{},
		&models.Customer{},
		&models.Staff{},
		&models.Owner{},
		&models.Booking{},
		&models.Payment{},
		&models.Chat{},
		&models.ChatMessage{},
		&models.Review{},
		&models.RevenueReport{},
	)

	if errMigrate != nil {
		panic("Gagal melakukan migrasi database: " + errMigrate.Error())
	}

	log.Println("Database berhasil terkoneksi dan tabel sudah ter-update!")
	DB = db
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
