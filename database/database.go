package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDB() {
	// Mengambil konfigurasi dari Environment Variables Docker
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Fallback jika tidak menemukan environment variables
	if host == "" {
		host = "localhost"
		port = "3306"
		user = "fleetify_user"
		password = "fleetifypassword"
		dbname = "fleetify_db"
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, host, port, dbname)

	database, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terkoneksi ke database! \n", err)
	}

	log.Println("Koneksi ke database MySQL via GORM berhasil!")
	DB = database
}