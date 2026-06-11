package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func init() {
	// Look for .env in the project root — services run from different subdirectories
	candidates := []string{".env", "../.env", "../../.env"}
	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			if err := godotenv.Load(p); err == nil {
				abs, _ := filepath.Abs(p)
				log.Println("Loaded env from", abs)
				return
			}
		}
	}
}

func GetEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func ConnectDB(schema string, models ...interface{}) *gorm.DB {
	host := GetEnv("DB_HOST", "localhost")
	user := GetEnv("DB_USER", "postgres")
	password := GetEnv("DB_PASSWORD", "")
	dbname := GetEnv("DB_NAME", "db_hotel_arca")
	port := GetEnv("DB_PORT", "5432")

	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s search_path=%s sslmode=disable",
		host, user, password, dbname, port, schema,
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("Gagal connect database: " + err.Error())
	}

	// Create schema if it doesn't exist (search_path alone won't do this)
	db.Exec(fmt.Sprintf("CREATE SCHEMA IF NOT EXISTS %s", schema))

	if len(models) > 0 {
		if err := db.AutoMigrate(models...); err != nil {
			panic("Gagal melakukan migrasi database: " + err.Error())
		}
		log.Println("Database berhasil terkoneksi dan tabel sudah ter-update!")
	} else {
		log.Println("Database berhasil terkoneksi!")
	}

	return db
}
