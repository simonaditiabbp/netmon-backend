package repository

import (
	"github.com/simonaditiabbp/netmon-backend/internal/domain"
	"gorm.io/gorm"
)

type DeviceTypeMapRepository struct {
	DB *gorm.DB
}

func NewDeviceTypeMapRepository(db *gorm.DB) *DeviceTypeMapRepository {
	return &DeviceTypeMapRepository{DB: db}
}

func (r *DeviceTypeMapRepository) AddDeviceTypes(deviceID uint, typeIDs []uint) error {
	maps := make([]domain.DeviceTypeMap, 0, len(typeIDs))
	for _, tid := range typeIDs {
		maps = append(maps, domain.DeviceTypeMap{DeviceID: deviceID, TypeID: tid})
	}
	return r.DB.Create(&maps).Error
}

func (r *DeviceTypeMapRepository) UpdateDeviceTypes(deviceID uint, typeIDs []uint) error {
	// Remove old mappings
	if err := r.DB.Where("device_id = ?", deviceID).Delete(&domain.DeviceTypeMap{}).Error; err != nil {
		return err
	}
	// Add new mappings
	return r.AddDeviceTypes(deviceID, typeIDs)
}

func (r *DeviceTypeMapRepository) GetDeviceTypes(deviceID uint) ([]domain.DeviceType, error) {
	var types []domain.DeviceType
	if err := r.DB.Table("device_type_maps").
		Select("devices_types.*").
		Joins("join devices_types on device_type_maps.type_id = devices_types.id").
		Where("device_type_maps.device_id = ?", deviceID).
		Scan(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}
