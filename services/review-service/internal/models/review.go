package models

import (
	"time"
)

// Review is a struct desribing an Review
type Review struct {
    ID               uint           `gorm:"primaryKey" json:"id"`
    EventID          uint           `gorm:"not null;index:idx_event_user,unique" json:"event_id"`
    UserID           uint           `gorm:"not null;index:idx_event_user,unique" json:"user_id"`
    ReviewTime       time.Time      `gorm:"autoCreateTime" json:"Review_time"`           
    UpdatedAt        time.Time      `gorm:"autoUpdateTime" json:"updated_at"`                   
    Comment          string         `gorm:"type:text" json:"comment,omitempty"`                 
}


// JSON EXAMPLE

// {
//   "event_id": 123,
//   "user_id": 456,
//   "comment": "that event was fun"
// }

