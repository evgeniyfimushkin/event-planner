package models

import (
	"time"
)

// Event is a struct desribing an event 
type Event struct {
    ID              uint      `gorm:"primaryKey" json:"id"`
    Name            string    `gorm:"type:varchar(255);not null;index" json:"name"`
    Description     string    `gorm:"type:text" json:"description"`
    Category        string    `gorm:"type:varchar(100);index" json:"category"`
    MaxParticipants int       `gorm:"default:100;check:max_participants >= 1" json:"max_participants"`
    // base64 image
    ImageData       []byte    `gorm:"type:bytea" json:"image_data"`

    City            string    `gorm:"type:varchar(100);not null;index" json:"city"`
    Address         string    `gorm:"type:varchar(255)" json:"address"`
    Latitude        float64   `gorm:"type:double precision;check:latitude >= -90 AND latitude <= 90" json:"latitude"`
    Longitude       float64   `gorm:"type:double precision;check:longitude >= -180 AND longitude <= 180" json:"longitude"`

    StartTime       time.Time `gorm:"not null;index" json:"start_time"`
    EndTime         time.Time `gorm:"not null;index" json:"end_time"`
    Status          string    `gorm:"type:varchar(50);not null;default:'active'" json:"status"`

    CreatedBy       string    `gorm:"not null;index" json:"created_by"`
    CreatedAt       time.Time `gorm:"autoCreateTime;default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt       time.Time `gorm:"autoUpdateTime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// JSON EXAMPLE

// {
//   "name": "Баскетбол",
//   "description": "Играем баскет на улице",
//   "category": "Спорт",
//   "max_participants": 30,
//   "image_data": "iVBORw0KGgoAAAANSUhEUgAAA...", 
//   "city": "Новосибирск",
//   "address": "Карла Маркса 37",
//   "latitude": 54.989688,
//   "longitude": 82.902014,
//   "start_time": "2025-05-15T17:00:00+07:00",
//   "end_time": "2025-05-15T23:00:00+07:00",
//   "status": "active",
//   "created_by": "evgeniyfimushkin"
// }

