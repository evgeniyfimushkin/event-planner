package models

import "time"

type Event struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    Name            string    `gorm:"type:varchar(255);not null" json:"name"`
    Description     string    `gorm:"type:text" json:"description"`
    Category        string    `gorm:"type:varchar(100)" json:"category"`
    ImagePath       string    `gorm:"type:varchar(255)" json:"image_path"`
    Location        string    `gorm:"type:varchar(255)" json:"location"`
    StartTime       time.Time `gorm:"not null" json:"start_time"`
    EndTime         time.Time `gorm:"not null" json:"end_time"`
    MaxParticipants int       `gorm:"default:100" json:"max_participants"`
    CreatedBy       uint      `gorm:"not null" json:"created_by"`
    CreatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt       time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
    Status          string    `gorm:"type:varchar(50);default:'active'" json:"status"`
}

