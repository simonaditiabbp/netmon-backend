package domain

import "time"

type Log struct {
	ID        uint      `gorm:"primaryKey"`
	DeviceID  uint      `gorm:"not null"`
	OldStatus string    `gorm:"column:oldstatus;not null"`
	NewStatus string    `gorm:"column:newstatus;not null"`
	Logtime   time.Time `gorm:"column:log_time;autoCreateTime"`
}

// TableName overrides the default table name
func (Log) TableName() string {
	return "logss"
}
