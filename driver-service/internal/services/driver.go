package services

import (
	"context"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
)

type DriverStorage interface {
	GetByID(ctx context.Context, driverID string) (*models.Driver, error)
	Update(ctx context.Context, driverID string, driverUpdate models.Driver) error
	ChangeStatus(ctx context.Context, driverID string, status models.DriverStatus) error
	SetInactive(ctx context.Context, driverID string) error
}

type RatingCounter interface {
	GetRating(ctx context.Context, driverID string) (float32, error)
}

type DriverService struct {
	driverStorage DriverStorage
	ratingCounter RatingCounter
}

func NewDriverService(driverStorage DriverStorage, ratingCounter RatingCounter) *DriverService {
	return &DriverService{
		driverStorage: driverStorage,
		ratingCounter: ratingCounter,
	}
}

func (s *DriverService) GetProfile(ctx context.Context, driverID string) (*models.Driver, error) {
	driver, err := s.driverStorage.GetByID(ctx, driverID)
	if err != nil {
		return &models.Driver{}, err
	}

	rating, err := s.ratingCounter.GetRating(ctx, driverID)
	if err != nil {
		return &models.Driver{}, err
	}

	driver.Rating = rating

	return driver, nil
}

func (s *DriverService) Update(ctx context.Context, driverID string, driverUpdate models.Driver) error {
	driver, err := s.driverStorage.GetByID(ctx, driverID)
	if err != nil {
		return err
	}

	if driverUpdate.Name != "" {
		driver.Name = driverUpdate.Name
	}

	if driverUpdate.Surname != "" {
		driver.Surname = driverUpdate.Surname
	}

	if driverUpdate.Phone != "" {
		driver.Phone = driverUpdate.Phone
	}

	if driverUpdate.Email != "" {
		driver.Email = driverUpdate.Email
	}

	return s.driverStorage.Update(ctx, driverID, *driver)
}

func (s *DriverService) ChangeStatus(ctx context.Context, driverID string, status models.DriverStatus) error {
	return s.driverStorage.ChangeStatus(ctx, driverID, status)
}

func (s *DriverService) SetInactive(ctx context.Context, driverID string) error {
	return s.driverStorage.SetInactive(ctx, driverID)
}
