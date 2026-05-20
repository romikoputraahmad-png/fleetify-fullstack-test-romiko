package main

import (
	"fleetify/database"
	"fleetify/handlers"
	"fleetify/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	database.ConnectDB()

	app := fiber.New()

	// Endpoint Testing
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Aplikasi Fleetify berjalan mantap!")
	})

	// --- SETUP ROUTES ---
	// --- SETUP ROUTES ---
	api := app.Group("/api")

	// F-04: Riwayat laporan (Bisa diakses SA dan APPROVAL)
	api.Get("/reports", middleware.RoleCheck("SA", "APPROVAL"), handlers.GetAllReports)

	// F-01: Membuat laporan baru (HANYA SA)
	api.Post("/reports", middleware.RoleCheck("SA"), handlers.CreateReport)

	// F-02: Menyetujui laporan (HANYA APPROVAL)
	api.Patch("/reports/:id/approve", middleware.RoleCheck("APPROVAL"), handlers.ApproveReport)

	// F-03: Menyelesaikan laporan (HANYA SA)
	api.Patch("/reports/:id/complete", middleware.RoleCheck("SA"), handlers.CompleteReport)

	log.Println("Server berjalan di port 3000...")
	log.Fatal(app.Listen(":3000"))
}
