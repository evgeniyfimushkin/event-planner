package models

import (
	"time"
)

// Registration is a struct desribing an registration
type Registration struct {
    ID               uint           `gorm:"primaryKey" json:"id"`
    EventID          uint           `gorm:"not null;index:idx_event_user,unique" json:"event_id"`
    UserID           uint           `gorm:"not null;index:idx_event_user,unique" json:"user_id"`
    RegistrationTime time.Time      `gorm:"autoCreateTime" json:"registration_time"`           
    Status           string         `gorm:"type:varchar(50);not null;default:'registered'" json:"status"`
    UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                   
    Comment          string         `gorm:"type:text" json:"comment,omitempty"`                 
}


// JSON EXAMPLE

// {
//   "event_id": 123,
//   "user_id": 456,
//   "comment": "ку-ку ёпта"
// }

