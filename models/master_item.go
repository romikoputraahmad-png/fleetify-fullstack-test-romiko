package models

import "time"

type MasterItem struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    ItemName  string    `gorm:"type:varchar(100);not null" json:"item_name"`
    Type      string    `gorm:"type:enum('PART', 'SERVICE');not null" json:"type"`
    Price     float64   `gorm:"type:decimal(10,2);not null" json:"price"`
    CreatedAt time.Time `json:"created_at"`
}