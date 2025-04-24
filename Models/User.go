package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name         string `gorm:"not null"`
	ProfileImage string
	Role         string
	Course       string
	Year         string
	Password     string `gorm:"not null"`
	Email        string `gorm:"unique;not null"`
	Posts        []Post `gorm:"foreignKey:UserID"`
}
