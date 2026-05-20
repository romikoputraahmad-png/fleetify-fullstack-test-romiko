package handlers

import (
	"fleetify/models"
	"fleetify/repositories"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Struct untuk menerima data JSON dari Frontend
type CreateReportRequest struct {
	VehicleID    uint                `json:"vehicle_id"`
	Odometer     int                 `json:"odometer"`
	Complaint    string              `json:"complaint"`
	InitialPhoto string              `json:"initial_photo"`
	Items        []models.ReportItem `json:"items"`
}

func CreateReport(c *fiber.Ctx) error {
	// Ambil data user yang sedang login (dari Middleware)
	user := c.Locals("user").(models.User)

	// Parsing JSON dari body request
	var req CreateReportRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format data tidak valid",
		})
	}

	// Buat objek Laporan Header
	report := models.MaintenanceReport{
		VehicleID:    req.VehicleID,
		CreatedBy:    user.ID, // Diambil otomatis dari user yang request
		Odometer:     req.Odometer,
		Complaint:    req.Complaint,
		InitialPhoto: req.InitialPhoto,
		Status:       "PENDING_APPROVAL", // Sesuai requirement soal
	}

	// Panggil fungsi Transaksi dari Repository
	if err := repositories.CreateReport(&report, req.Items); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan laporan: " + err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":   "Laporan berhasil dibuat",
		"report_id": report.ID,
	})
}
func ApproveReport(c *fiber.Ctx) error {
	// Fitur F-02: Approval menyetujui laporan
	reportID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID Laporan tidak valid"})
	}

	if err := repositories.UpdateReportStatus(uint(reportID), "APPROVED", nil); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyetujui laporan"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Laporan berhasil disetujui (APPROVED)"})
}

// Struct untuk menerima payload foto bukti dari Frontend
type CompleteReportRequest struct {
	ProofPhoto string `json:"proof_photo"`
}

func CompleteReport(c *fiber.Ctx) error {
	// Fitur F-03: SA menyelesaikan laporan dengan bukti foto
	reportID, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "ID Laporan tidak valid"})
	}

	var req CompleteReportRequest
	if err := c.BodyParser(&req); err != nil || req.ProofPhoto == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Foto bukti pengerjaan wajib diisi"})
	}

	if err := repositories.UpdateReportStatus(uint(reportID), "COMPLETED", &req.ProofPhoto); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyelesaikan laporan"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Laporan berhasil diselesaikan (COMPLETED)"})
}

func GetAllReports(c *fiber.Ctx) error {
	// Fitur F-04: Menampilkan riwayat laporan
	reports, err := repositories.GetAllReports()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data laporan"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": reports})
}
