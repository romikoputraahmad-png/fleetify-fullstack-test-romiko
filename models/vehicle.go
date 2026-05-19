package models

import "time"

type Vehicle struct {
    ID           uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    LicensePlate string    `gorm:"type:varchar(20);not null;unique" json:"license_plate"`
    Model        string    `gorm:"type:varchar(100);not null" json:"model"`
    CreatedAt    time.Time `json:"created_at"`
}