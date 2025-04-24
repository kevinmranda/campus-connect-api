package models

import "gorm.io/gorm"

type Job struct {
	gorm.Model
	Title string `gorm:"not null"`
	Description string `gorm:"not null"`
	Company string `gorm:"not null"`
	Link string `gorm:"not null"`
}