package models

import "time"

type Location struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    City      string    `gorm:"type:varchar(100);not null" json:"city"`
    Address   string    `gorm:"type:varchar(255)" json:"address"`
    Latitude  float64   `gorm:"type:double" json:"latitude"`
    Longitude float64   `gorm:"type:double" json:"longitude"`
    CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

