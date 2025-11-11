package domain

import "time"

type DeviceType struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	TypeName    string `gorm:"unique;not null;column:type_name" json:"type_name"`
	Icon        string
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	CreatedBy   string    `json:"created_by"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	UpdatedBy   string    `json:"updated_by"`
}

type DeviceTypeMap struct {
	ID       uint `gorm:"primaryKey;autoIncrement"`
	DeviceID uint `gorm:"not null;index"`
	TypeID   uint `gorm:"not null;index"`
}

func (DeviceType) TableName() string {
	return "devices_types"
}
