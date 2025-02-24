package models

type User struct {
    ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"not null" json:"username"`
	Email    string `gorm:"not null; unique" json:"email"`
	PassHash string `gorm:"not null" json:"passhash"`
}

