package usecase

import (
	"github.com/simonaditiabbp/netmon-backend/internal/domain"
	"github.com/simonaditiabbp/netmon-backend/internal/repository"
)

type LocationUsecase struct {
	Repo *repository.LocationRepository
}

func NewLocationUsecase(repo *repository.LocationRepository) *LocationUsecase {
	return &LocationUsecase{Repo: repo}
}

func (u *LocationUsecase) CreateLocation(loc *domain.Location) error {
	return u.Repo.CreateLocation(loc)
}

func (u *LocationUsecase) UpdateLocation(loc *domain.Location) error {
	return u.Repo.UpdateLocation(loc)
}

func (u *LocationUsecase) GetAllLocations() ([]domain.Location, error) {
	return u.Repo.GetAllLocations()
}

func (u *LocationUsecase) GetLocationByID(id uint) (*domain.Location, error) {
	return u.Repo.GetLocationByID(id)
}

func (u *LocationUsecase) DeleteLocation(id uint) error {
	return u.Repo.DeleteLocation(id)
}
