package models

type User struct {
    ID       uint   `gorm:"primaryKey" json:"id"`
	Username string `gorm:"unique;not null" json:"username"`
	Email    string `gorm:"unique" json:"email"`
	PassHash string `gorm:"not null" json:"passhash"`
    Role     string `gorm:"not null" json:"role"`
}

