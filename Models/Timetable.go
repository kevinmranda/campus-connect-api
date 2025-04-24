package models

import (
	"gorm.io/gorm"
	"time"
)

type Timetable struct {
	gorm.Model
	Day         string    `gorm:"not null"`
	Subject     string    `gorm:"not null"`
	SubjectCode string    `gorm:"not null"`
	Faculty     string    `gorm:"not null"`
	Room        string    `gorm:"not null"`
	Time        time.Time `gorm:"not null"`
	Instructor  string    `gorm:"not null"`
}
