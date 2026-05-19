package main

import (
	"fleetify/database" // Sesuaikan dengan nama module-mu jika berbeda
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// Panggil koneksi database
	database.ConnectDB()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Aplikasi Fleetify berjalan mantap dan terhubung ke DB!")
	})

	log.Println("Server berjalan di port 3000...")
	log.Fatal(app.Listen(":3000"))
}