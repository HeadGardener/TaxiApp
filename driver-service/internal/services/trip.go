package services

import (
	"context"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
)

type TripStorage interface {
	GetAll(ctx context.Context, driverID string) ([]*models.Trip, error)
}

type TripService struct {
	tripStorage TripStorage
}

func NewTripService(tripStorage TripStorage) *TripService {
	return &TripService{tripStorage: tripStorage}
}

func (s *TripService) ViewAll(ctx context.Context, driverID string) ([]*models.Trip, error) {
	return s.tripStorage.GetAll(ctx, driverID)
}
