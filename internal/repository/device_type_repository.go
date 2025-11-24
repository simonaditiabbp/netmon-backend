package repository

import (
	"github.com/simonaditiabbp/netmon-backend/internal/domain"
	"gorm.io/gorm"
)

type DeviceTypeRepository struct {
	DB *gorm.DB
}

func NewDeviceTypeRepository(db *gorm.DB) *DeviceTypeRepository {
	return &DeviceTypeRepository{DB: db}
}

func (r *DeviceTypeRepository) CreateDeviceType(dt *domain.DeviceType) error {
	return r.DB.Create(dt).Error
}

func (r *DeviceTypeRepository) UpdateDeviceType(dt *domain.DeviceType) error {
	return r.DB.Model(&domain.DeviceType{}).Where("id = ?", dt.ID).Updates(dt).Error
}

func (r *DeviceTypeRepository) GetAllDeviceTypes() ([]domain.DeviceType, error) {
	var types []domain.DeviceType
	if err := r.DB.Order("type_name ASC").Find(&types).Error; err != nil {
		return nil, err
	}
	return types, nil
}

func (r *DeviceTypeRepository) GetDeviceTypeByID(id uint) (*domain.DeviceType, error) {
	var dt domain.DeviceType
	if err := r.DB.First(&dt, id).Error; err != nil {
		return nil, err
	}
	return &dt, nil
}
