package storage

import (
	"context"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
	"github.com/jmoiron/sqlx"
)

type DriverStorage struct {
	db *sqlx.DB
}

func NewDriverStorage(db *sqlx.DB) *DriverStorage {
	return &DriverStorage{db: db}
}

func (s *DriverStorage) Create(ctx context.Context, driver *models.Driver) (string, error) {
	if _, err := s.db.ExecContext(ctx, `INSERT INTO drivers
    (id, name, surname, phone, email, taxi_type, balance, password_hash, rating, driver_status, registration, is_active)
    VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`,
		driver.ID,
		driver.Name,
		driver.Surname,
		driver.Phone,
		driver.Email,
		driver.TaxiType,
		driver.Balance,
		driver.Password,
		driver.Rating,
		driver.DriverStatus,
		driver.Registration,
		driver.IsActive); err != nil {
		return "", err
	}

	return driver.ID, nil
}

func (s *DriverStorage) GetByPhone(ctx context.Context, phone string) (*models.Driver, error) {
	var driver models.Driver

	if err := s.db.SelectContext(ctx, &driver, `SELECT * FROM drivers WHERE phone=$1`, phone); err != nil {
		return nil, err
	}

	return &driver, nil
}

func (s *DriverStorage) GetByID(ctx context.Context, driverID string) (*models.Driver, error) {
	var driver models.Driver

	if err := s.db.SelectContext(ctx, &driver, `SELECT * FROM drivers WHERE id=$1`, driverID); err != nil {
		return nil, err
	}

	return &driver, nil
}

func (s *DriverStorage) Update(ctx context.Context, driverID string, driverUpdate models.Driver) error {
	if _, err := s.db.ExecContext(ctx,
		`UPDATE drivers SET name=$1, surname=$2, phone=$3, email=$4 WHERE id=$5`,
		driverUpdate.Name,
		driverUpdate.Surname,
		driverUpdate.Phone,
		driverUpdate.Email,
		driverID); err != nil {
		return err
	}

	return nil
}

func (s *DriverStorage) ChangeStatus(ctx context.Context, driverID string, status models.DriverStatus) error {
	if _, err := s.db.ExecContext(ctx,
		`UPDATE drivers SET status=$1 WHERE id=$2`,
		status,
		driverID); err != nil {
		return err
	}

	return nil
}

func (s *DriverStorage) SetInactive(ctx context.Context, driverID string) error {
	if _, err := s.db.ExecContext(ctx,
		`UPDATE drivers SET is_active=false WHERE id=$1`,
		driverID); err != nil {
		return err
	}

	return nil
}
