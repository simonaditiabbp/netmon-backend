package usecase

import (
	"github.com/simonaditiabbp/netmon-backend/internal/domain"
	"github.com/simonaditiabbp/netmon-backend/internal/repository"
)

type DeviceTypeMapUsecase struct {
	Repo *repository.DeviceTypeMapRepository
}

func NewDeviceTypeMapUsecase(repo *repository.DeviceTypeMapRepository) *DeviceTypeMapUsecase {
	return &DeviceTypeMapUsecase{Repo: repo}
}

func (u *DeviceTypeMapUsecase) AddDeviceTypes(deviceID uint, typeIDs []uint) error {
	return u.Repo.AddDeviceTypes(deviceID, typeIDs)
}

func (u *DeviceTypeMapUsecase) UpdateDeviceTypes(deviceID uint, typeIDs []uint) error {
	return u.Repo.UpdateDeviceTypes(deviceID, typeIDs)
}

func (u *DeviceTypeMapUsecase) GetDeviceTypes(deviceID uint) ([]domain.DeviceType, error) {
	return u.Repo.GetDeviceTypes(deviceID)
}
