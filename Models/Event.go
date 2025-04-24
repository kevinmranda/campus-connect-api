package models

import (
	"gorm.io/gorm"
	"time"
)

type Event struct {
	gorm.Model
	Quarter      string    `gorm:"not null"`
	Month        string    `gorm:"not null"`
	Date         time.Time `gorm:"not null"`
	Title        string    `gorm:"not null"`
	Participants string
}