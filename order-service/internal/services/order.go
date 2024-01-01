package services

import (
	"context"
	"errors"

	"github.com/HeadGardener/TaxiApp/order-service/internal/models"
	"github.com/google/uuid"
)

var (
	ErrNotDriversOrder = errors.New("this order not of this driver")
)

type UserServiceGRPCClient interface {
	AcceptOrder(ctx context.Context, userID, orderID, driverID string, status int32) error
	CompleteOrder(ctx context.Context, userID, orderID string, status int32) error
}

type OrderStorage interface {
	GetAll(ctx context.Context) ([]models.Order, error)
	Save(ctx context.Context, order *models.Order) (string, error)
	GetByID(ctx context.Context, orderID string) (models.Order, error)
	AddComment(ctx context.Context, orderID, comment string) error
	UpdateStatus(ctx context.Context, orderID string, status int32) error
	UpdateRating(ctx context.Context, orderID string, rating float64) error
}

type OrderService struct {
	client       UserServiceGRPCClient
	orderStorage OrderStorage
}

func NewOrderService(client UserServiceGRPCClient, orderStorage OrderStorage) *OrderService {
	return &OrderService{
		client:       client,
		orderStorage: orderStorage,
	}
}

func (s *OrderService) GetAll(ctx context.Context) ([]models.Order, error) {
	return s.orderStorage.GetAll(ctx)
}

func (s *OrderService) CreateOrder(ctx context.Context, order *models.Order) (string, error) {
	order.ID = uuid.NewString()
	order.Status = models.Creating

	return s.orderStorage.Save(ctx, order)
}

func (s *OrderService) AddComment(ctx context.Context, orderID, comment string) error {
	return s.orderStorage.AddComment(ctx, orderID, comment)
}

func (s *OrderService) ProcessOrder(ctx context.Context, orderID, driverID string, status int32) error {
	order, err := s.orderStorage.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.DriverID != driverID {
		return ErrNotDriversOrder
	}

	if err = s.client.AcceptOrder(ctx, order.UserID, orderID, driverID, status); err != nil {
		return err
	}

	return s.orderStorage.UpdateStatus(ctx, orderID, status)
}

func (s *OrderService) CompleteOrder(ctx context.Context, driverID, orderID string,
	status int32, rating float64) error {
	order, err := s.orderStorage.GetByID(ctx, orderID)
	if err != nil {
		return err
	}

	if order.DriverID != driverID {
		return ErrNotDriversOrder
	}

	if err = s.client.CompleteOrder(ctx, order.UserID, orderID, status); err != nil {
		return err
	}

	return s.orderStorage.UpdateRating(ctx, orderID, rating)
}
