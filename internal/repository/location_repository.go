package repository

import (
	"github.com/simonaditiabbp/netmon-backend/internal/domain"
	"gorm.io/gorm"
)

type LocationRepository struct {
	DB *gorm.DB
}

func NewLocationRepository(db *gorm.DB) *LocationRepository {
	return &LocationRepository{DB: db}
}

func (r *LocationRepository) CreateLocation(loc *domain.Location) error {
	return r.DB.Create(loc).Error
}

func (r *LocationRepository) UpdateLocation(loc *domain.Location) error {
	return r.DB.Model(&domain.Location{}).Where("id = ?", loc.ID).Updates(loc).Error
}

func (r *LocationRepository) GetAllLocations() ([]domain.Location, error) {
	var locs []domain.Location
	if err := r.DB.Find(&locs).Error; err != nil {
		return nil, err
	}
	return locs, nil
}

func (r *LocationRepository) GetLocationByID(id uint) (*domain.Location, error) {
	var loc domain.Location
	if err := r.DB.First(&loc, id).Error; err != nil {
		return nil, err
	}
	return &loc, nil
}

func (r *LocationRepository) DeleteLocation(id uint) error {
	return r.DB.Delete(&domain.Location{}, id).Error
}
