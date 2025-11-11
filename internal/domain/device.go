package domain

import "time"

type Device struct {
	ID         uint   `gorm:"primaryKey;autoIncrement"`
	Name       string `gorm:"not null"`
	IP         string
	URL        string
	Status     string
	LastOnline time.Time `gorm:"column:lastonline"`
	Icon       string
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	CreatedBy  string
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
	UpdatedBy  string
}
