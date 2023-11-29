package storage

import (
	"context"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
	"github.com/jmoiron/sqlx"
)

type TripStorage struct {
	db *sqlx.DB
}

func NewTripStorage(db *sqlx.DB) *TripStorage {
	return &TripStorage{db: db}
}

func (s *TripStorage) GetAll(ctx context.Context, driverID string) ([]*models.Trip, error) {
	var trips []*models.Trip

	if err := s.db.SelectContext(ctx, &trips, `SELECT * FROM trips WHERE driver_id=$1`, driverID); err != nil {
		return nil, err
	}

	return trips, nil
}
