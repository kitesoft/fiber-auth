package model

import "github.com/jinzhu/gorm"

// User struct
type User struct {
	gorm.Model
	Username string `gorm:"unique;not null" json:"username"`
	Email    string `gorm:"unique;not null" json:"email"`
	Password string `json:"password"`
}
