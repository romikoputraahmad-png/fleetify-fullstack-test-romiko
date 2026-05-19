package models

import "time"

type User struct {
    ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
    Username  string    `gorm:"type:varchar(50);not null;unique" json:"username"`
    Role      string    `gorm:"type:enum('SA', 'APPROVAL');not null" json:"role"`
    CreatedAt time.Time `json:"created_at"`
}