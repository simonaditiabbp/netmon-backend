package usecase

import (
	"github.com/simonaditiabbp/netmon-backend/internal/domain"
	"github.com/simonaditiabbp/netmon-backend/internal/repository"
)

type DeviceTypeUsecase struct {
	Repo *repository.DeviceTypeRepository
}

func NewDeviceTypeUsecase(repo *repository.DeviceTypeRepository) *DeviceTypeUsecase {
	return &DeviceTypeUsecase{Repo: repo}
}

func (u *DeviceTypeUsecase) CreateDeviceType(dt *domain.DeviceType) error {
	return u.Repo.CreateDeviceType(dt)
}

func (u *DeviceTypeUsecase) UpdateDeviceType(dt *domain.DeviceType) error {
	return u.Repo.UpdateDeviceType(dt)
}

func (u *DeviceTypeUsecase) GetAllDeviceTypes() ([]domain.DeviceType, error) {
	return u.Repo.GetAllDeviceTypes()
}

func (u *DeviceTypeUsecase) GetDeviceTypeByID(id uint) (*domain.DeviceType, error) {
	return u.Repo.GetDeviceTypeByID(id)
}

func (u *DeviceTypeUsecase) CountDeviceTypeUsage(typeID uint) (int64, error) {
	var count int64
	err := u.Repo.DB.Model(&domain.DeviceTypeMap{}).Where("type_id = ?", typeID).Count(&count).Error
	return count, err
}

func (u *DeviceTypeUsecase) DeleteDeviceType(typeID uint) error {
	return u.Repo.DB.Delete(&domain.DeviceType{}, typeID).Error
}
