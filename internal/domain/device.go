package domain

import "time"

type Device struct {
	ID         uint   `gorm:"primaryKey"`
	Name       string `gorm:"not null"`
	IP         string
	URL        string
	Status     string
	LastOnline time.Time `gorm:"column:lastonline"`
	Icon       string
}
