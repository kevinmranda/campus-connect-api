package models

import (
	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Image       string
	Title       string `gorm:"not null"`
	Description string
	UserID      uint
	User        User
}
