package storage

import (
	"context"
	"fmt"
	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
	"github.com/jmoiron/sqlx"
)

type TripStorage struct {
	db *sqlx.DB
}

func NewTripStorage(db *sqlx.DB) *TripStorage {
	return &TripStorage{db: db}
}

func (s *TripStorage) Save(ctx context.Context, userID string, trip *models.Trip) error {
	var saveTripQuery = fmt.Sprintf(`INSERT INTO %s (id, user_id, taxi_type, "from", "to")
										VALUES($1,$2,$3,$4,$5)`, tripsTable)

	if _, err := s.db.ExecContext(ctx,
		saveTripQuery,
		trip.ID,
		userID,
		trip.TaxiType,
		trip.From,
		trip.To); err != nil {
		return err
	}

	return nil
}

func (s *TripStorage) GetAll(ctx context.Context, userID string) ([]models.Trip, error) {
	var getAllTripsQuery = fmt.Sprintf(`SELECT (taxi_type, driver, "from", "to", rating, date) FROM %s 
                                        WHERE user_id=$1`, tripsTable)

	var trips []models.Trip

	if err := s.db.SelectContext(ctx, &trips, getAllTripsQuery, userID); err != nil {
		return nil, err
	}

	return trips, nil
}
