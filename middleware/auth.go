package middleware

import (
	"fleetify/database"
	"fleetify/models"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

// RoleCheck memastikan pengguna memiliki role yang sesuai
func RoleCheck(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil X-User-ID dari Header
		userIDStr := c.Get("X-User-ID")
		if userIDStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Header X-User-ID tidak ditemukan",
			})
		}

		// Konversi string ke integer
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Format X-User-ID tidak valid",
			})
		}

		// 2. Cari User di Database
		var user models.User
		if err := database.DB.First(&user, userID).Error; err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "User tidak ditemukan",
			})
		}

		// 3. Cek apakah role user ada di dalam daftar allowedRoles
		roleAllowed := false
		for _, role := range allowedRoles {
			if user.Role == role {
				roleAllowed = true
				break
			}
		}

		if !roleAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Anda tidak memiliki akses (Role tidak sesuai)",
			})
		}

		// 4. Simpan data user ke context agar bisa dipakai di Handler berikutnya
		c.Locals("user", user)

		return c.Next()
	}
}