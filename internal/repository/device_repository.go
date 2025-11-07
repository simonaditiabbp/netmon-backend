package repository

import (
	"github.com/simonaditiabbp/netmon-backend/internal/domain"

	"gorm.io/gorm"
)

type DeviceRepository struct {
	DB *gorm.DB
}

func NewDeviceRepository(db *gorm.DB) *DeviceRepository {
	return &DeviceRepository{DB: db}
}

func (r *DeviceRepository) GetAllDevices() ([]domain.Device, error) {
	var devices []domain.Device
	if err := r.DB.Find(&devices).Error; err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *DeviceRepository) CreateDevice(device *domain.Device) error {
	return r.DB.Create(device).Error
}

func (r *DeviceRepository) CreateLog(log *domain.Log) error {
	return r.DB.Create(log).Error
}

func (r *DeviceRepository) UpdateDevice(device *domain.Device) error {
	return r.DB.Model(&domain.Device{}).Where("id = ?", device.ID).Updates(device).Error
}
