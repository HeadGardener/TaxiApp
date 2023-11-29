package services

import (
	"context"
	"errors"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
	"log/slog"
)

var (
	ErrNotEnoughBalance = errors.New("not enough balance")
)

type BalanceProcessor interface {
	Add(ctx context.Context, driverID string, money float64) error
	Withdraw(ctx context.Context, driverID string, money float64) error
}

type DriverProvider interface {
	GetByID(ctx context.Context, driverID string) (*models.Driver, error)
}

type BalanceService struct {
	balanceProcessor BalanceProcessor
	driverProvider   DriverProvider
}

func NewBalanceService(balanceProcessor BalanceProcessor, driverProvider DriverProvider) *BalanceService {
	return &BalanceService{
		balanceProcessor: balanceProcessor,
		driverProvider:   driverProvider,
	}
}

func (s *BalanceService) Add(ctx context.Context, driverID string, money float64) error {
	return s.balanceProcessor.Add(ctx, driverID, money)
}

// CashOut - provides driver ability to cash out his earnings
// this method is placeholder
func (s *BalanceService) CashOut(ctx context.Context, driverID string, credentials models.Credentials) error {
	driver, err := s.driverProvider.GetByID(ctx, driverID)
	if err != nil {
		return err
	}

	if driver.Balance < credentials.Money {
		return ErrNotEnoughBalance
	}

	if err = s.balanceProcessor.Withdraw(ctx, driverID, credentials.Money); err != nil {
		return err
	}

	slog.Info("send %f on card %s", credentials.Money, credentials.Card)

	return nil
}
