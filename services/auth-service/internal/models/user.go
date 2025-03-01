package models

import "time"

type User struct {
    ID              int    `gorm:"primaryKey" json:"id"`
	Username        string `gorm:"unique;not null" json:"username"`
	Email           string `gorm:"unique" json:"email"`
	PassHash        string `gorm:"not null" json:"passhash"`
    Role            string `gorm:"not null" json:"role"`
    CreatedAt       time.Time `gorm:"autoCreateTime;default:CURRENT_TIMESTAMP" json:"created_at"`
    UpdatedAt       time.Time `gorm:"autoUpdateTime;default:CURRENT_TIMESTAMP" json:"updated_at"`
}

// JSON EXAMPLE

// {
//     "username": "ivan",
//     "passhash": "asdfhj87314gy8asdfh3478ysuadf",
//     "email": "ivan@gmail.com"
// }
