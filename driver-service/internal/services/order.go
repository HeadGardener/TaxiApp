package services

import (
	"context"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
	"time"
)

const (
	reqTimeout = 5 * time.Second
)

type GRPCClient interface {
}

type OrderService struct {
	client GRPCClient
}

func NewOrderService(client GRPCClient) *OrderService {
	return &OrderService{
		client: client,
	}
}

func (s *OrderService) StartAcceptOrders(ctx context.Context, driverID string, taxiType models.TaxiType) error {
	return nil
}

func (s *OrderService) Complete(ctx context.Context, driverID, orderID, status string) error {
	return nil
}
