package repositories

import (
	"bytes"
	"encoding/json"
	"errors"
	"fleetify/database"
	"fleetify/models"
	"log"
	"net/http"

	"gorm.io/gorm"
)

// CreateReport menyimpan laporan dan detail item dalam satu transaksi Atomic
func CreateReport(report *models.MaintenanceReport, items []models.ReportItem) error {
	// Memulai Transaksi GORM
	return database.DB.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan Header Laporan
		if err := tx.Create(report).Error; err != nil {
			return err
		}

		// 2. Simpan Detail Item
		for i := range items {
			var masterItem models.MasterItem
			// Ambil harga dari master item untuk di-snapshot
			if err := tx.First(&masterItem, items[i].ItemID).Error; err != nil {
				return errors.New("master item tidak ditemukan")
			}

			items[i].ReportID = report.ID
			items[i].PriceSnapshot = masterItem.Price

			if err := tx.Create(&items[i]).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// UpdateReportStatus mengubah status laporan dan menyimpan foto bukti jika ada
func UpdateReportStatus(reportID uint, status string, proofPhoto *string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if proofPhoto != nil {
		updates["proof_photo"] = *proofPhoto
	}

	err := database.DB.Model(&models.MaintenanceReport{}).Where("id = ?", reportID).Updates(updates).Error

	// Fitur Bonus B-02: Trigger Webhook dengan Goroutine (Asynchronous)
	if err == nil && (status == "APPROVED" || status == "COMPLETED") {
		go sendWebhook(reportID, status) // Menambahkan kata 'go' membuatnya berjalan asinkron
	}

	return err
}

// GetAllReports mengambil semua riwayat laporan beserta data kendaraan dan pembuatnya
func GetAllReports() ([]models.MaintenanceReport, error) {
	var reports []models.MaintenanceReport
	// GORM Preload digunakan untuk melakukan JOIN tabel secara otomatis
	err := database.DB.Preload("Vehicle").Preload("User").Order("created_at desc").Find(&reports).Error
	return reports, err
}

// sendWebhook mengirimkan HTTP POST secara asinkron
func sendWebhook(reportID uint, status string) {
	payload := map[string]interface{}{
		"report_id": reportID,
		"status":    status,
		"message":   "Status laporan pemeliharaan telah berubah",
	}
	jsonPayload, _ := json.Marshal(payload)

	// URL dummy untuk simulasi Webhook
	webhookURL := "https://webhook.site/dummy-url-fleetify"

	// Melakukan HTTP POST
	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		log.Println("[-] Gagal mengirim webhook:", err)
		return
	}
	defer resp.Body.Close()
	log.Printf("[+] Webhook berhasil dikirim untuk Report ID %d dengan status %s\n", reportID, status)
}
