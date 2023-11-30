package services

import (
	"context"
	"errors"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
)

var ErrInvalidDriverStatus = errors.New("invalid driver status")

type DriverStorage interface {
	GetByID(ctx context.Context, driverID string) (*models.Driver, error)
	Update(ctx context.Context, driverID string, driverUpdate models.Driver) error
	ChangeStatus(ctx context.Context, driverID string, status models.DriverStatus) error
	SetInactive(ctx context.Context, driverID string) error
}

type DriverService struct {
	driverStorage DriverStorage
}

func NewDriverService(driverStorage DriverStorage) *DriverService {
	return &DriverService{
		driverStorage: driverStorage,
	}
}

func (s *DriverService) GetProfile(ctx context.Context, driverID string) (*models.Driver, error) {
	driver, err := s.driverStorage.GetByID(ctx, driverID)
	if err != nil {
		return &models.Driver{}, err
	}

	driver.Password = ""

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

func (s *DriverService) ChangeStatus(ctx context.Context, driverID string, status string) error {
	st, ok := models.DriverStatusesStr[status]
	if !ok {
		return ErrInvalidDriverStatus
	}

	return s.driverStorage.ChangeStatus(ctx, driverID, st)
}

func (s *DriverService) SetInactive(ctx context.Context, driverID string) error {
	return s.driverStorage.SetInactive(ctx, driverID)
}
