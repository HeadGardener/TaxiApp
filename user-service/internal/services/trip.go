package services

import (
	"context"

	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
)

type TripProcessor interface {
	GetAll(ctx context.Context, userID string) ([]models.Trip, error)
}

type TripService struct {
	tripProcessor TripProcessor
}

func NewTripService(tripStorage TripProcessor) *TripService {
	return &TripService{tripProcessor: tripStorage}
}

func (s *TripService) ViewAll(ctx context.Context, userID string) ([]models.Trip, error) {
	return s.tripProcessor.GetAll(ctx, userID)
}
