package models

import "time"

type MaintenanceReport struct {
    ID           uint         `gorm:"primaryKey;autoIncrement" json:"id"`
    VehicleID    uint         `gorm:"not null" json:"vehicle_id"`
    CreatedBy    uint         `gorm:"not null" json:"created_by"`
    Odometer     int          `gorm:"not null" json:"odometer"`
    Complaint    string       `gorm:"type:text;not null" json:"complaint"`
    Status       string       `gorm:"type:enum('PENDING_APPROVAL', 'APPROVED', 'COMPLETED');default:'PENDING_APPROVAL'" json:"status"`
    InitialPhoto string       `gorm:"type:varchar(255)" json:"initial_photo"`
    ProofPhoto   string       `gorm:"type:varchar(255)" json:"proof_photo"`
    CreatedAt    time.Time    `json:"created_at"`

    // Relasi
    Vehicle      Vehicle      `gorm:"foreignKey:VehicleID" json:"vehicle"`
    User         User         `gorm:"foreignKey:CreatedBy" json:"user"`
    ReportItems  []ReportItem `gorm:"foreignKey:ReportID" json:"report_items"`
}

type ReportItem struct {
    ID            uint       `gorm:"primaryKey;autoIncrement" json:"id"`
    ReportID      uint       `gorm:"not null" json:"report_id"`
    ItemID        uint       `gorm:"not null" json:"item_id"`
    Quantity      int        `gorm:"not null" json:"quantity"`
    PriceSnapshot float64    `gorm:"type:decimal(10,2);not null" json:"price_snapshot"`

    // Relasi
    MasterItem    MasterItem `gorm:"foreignKey:ItemID" json:"master_item"`
}