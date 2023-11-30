package services

import (
	"context"
	"errors"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/lib/auth"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/lib/hash"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
	"github.com/google/uuid"
	"time"
)

var (
	ErrNotActive       = errors.New("unable to get access to this account")
	ErrInvalidPassword = errors.New("invalid password")
)

type DriverProcessor interface {
	Create(ctx context.Context, driver *models.Driver) (string, error)
	GetByPhone(ctx context.Context, phone string) (*models.Driver, error)
}

type TokenStorage interface {
	Add(ctx context.Context, driverID, token string) error
	Check(ctx context.Context, driverID, token string) error
	Delete(ctx context.Context, driverID string) error
}

type IdentityService struct {
	driverProcessor DriverProcessor
	tokenStorage    TokenStorage
}

func NewIdentityService(driverProcessor DriverProcessor, tokenStorage TokenStorage) *IdentityService {
	return &IdentityService{
		driverProcessor: driverProcessor,
		tokenStorage:    tokenStorage,
	}
}

func (s *IdentityService) SignUp(ctx context.Context, driver *models.Driver) (string, error) {
	driver.ID = uuid.NewString()
	driver.Balance = 0.0
	driver.Password = hash.GetPasswordHash(driver.Password)
	driver.Rating = 0.0
	driver.DriverStatus = models.Disable
	driver.Registration = time.Now()
	driver.IsActive = true

	return s.driverProcessor.Create(ctx, driver)
}

func (s *IdentityService) SignIn(ctx context.Context, phone, password string) (string, error) {
	driver, err := s.driverProcessor.GetByPhone(ctx, phone)
	if err != nil {
		return "", err
	}

	if !driver.IsActive {
		return "", ErrNotActive
	}

	if !hash.CheckPassword([]byte(driver.Password), password) {
		return "", ErrInvalidPassword
	}

	token, err := auth.GenerateToken(driver.ID, driver.Phone, driver.TaxiType)
	if err != nil {
		return "", err
	}

	if err = s.tokenStorage.Add(ctx, driver.ID, token); err != nil {
		return "", err
	}

	return token, nil
}

func (s *IdentityService) Check(ctx context.Context, driverID, token string) error {
	return s.tokenStorage.Check(ctx, driverID, token)
}

func (s *IdentityService) LogOut(ctx context.Context, driverID string) error {
	return s.tokenStorage.Delete(ctx, driverID)
}
