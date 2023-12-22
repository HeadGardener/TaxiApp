package services

import (
	"context"
	"time"

	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
)

const (
	reqTimeout = 5 * time.Second
)

type GRPCClient interface {
	AddDriver(ctx context.Context, driverID string, taxiType models.TaxiType) error
	ProcessOrder(ctx context.Context, driverID, orderID, status string) error
	CompleteOrder(ctx context.Context, driverID, orderID, status string, rating float64) error
}

type OrderStorage interface {
	Add(order *models.Order) error
	GetByDriverID(driverID string) (models.Order, error)
	Delete(driverID string) error
}

type OrderService struct {
	client       GRPCClient
	orderStorage OrderStorage
}

func NewOrderService(client GRPCClient, orderStorage OrderStorage) *OrderService {
	return &OrderService{
		client:       client,
		orderStorage: orderStorage,
	}
}

func (s *OrderService) Add(order *models.Order) error {
	return s.orderStorage.Add(order)
}

func (s *OrderService) GetInLine(ctx context.Context, driverID string, taxiType models.TaxiType) error {
	reqCtx, cancel := context.WithTimeout(ctx, reqTimeout)
	defer cancel()

	return s.client.AddDriver(reqCtx, driverID, taxiType)
}

func (s *OrderService) ProcessOrder(ctx context.Context, driverID, orderID, status string) error {
	reqCtx, cancel := context.WithTimeout(ctx, reqTimeout)
	defer cancel()

	return s.client.ProcessOrder(reqCtx, driverID, orderID, status)
}

func (s *OrderService) Complete(ctx context.Context, driverID, orderID, status string, rating float64) error {
	if err := s.orderStorage.Delete(driverID); err != nil {
		return err
	}

	reqCtx, cancel := context.WithTimeout(ctx, reqTimeout)
	defer cancel()

	return s.client.CompleteOrder(reqCtx, driverID, orderID, status, rating)
}

func (s *OrderService) CurrentOrder(driverID string) (models.Order, error) {
	return s.orderStorage.GetByDriverID(driverID)
}
