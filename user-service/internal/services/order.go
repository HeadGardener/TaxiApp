package services

import (
	"context"
	"errors"
	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
	"time"
)

var (
	ErrNotYourOrder = errors.New("this not your order")
)

const (
	reqTimeout = 5 * time.Second
)

type GRPCClient interface {
	CreateOrder(ctx context.Context, userID string, order *models.Order) (string, error)
	SendComment(ctx context.Context, orderID, comment string) error
}

type OrderStorage interface {
	Save(orderID string, order *models.UserOrder) error
	Get(orderID string) (*models.UserOrder, error)
	Update(orderID, status string) error
	Delete(orderID string) error
}

type TripSaver interface {
	Save(ctx context.Context, orderID, userID string, order *models.Order) error
	UpdateDriver(ctx context.Context, orderID, driverID string) error
}

type OrderService struct {
	client       GRPCClient
	orderStorage OrderStorage
	tripSaver    TripSaver
}

func NewOrderService(client GRPCClient, orderStorage OrderStorage, treipSaver TripSaver) *OrderService {
	return &OrderService{
		client:       client,
		orderStorage: orderStorage,
		tripSaver:    treipSaver,
	}
}

func (s *OrderService) Update(ctx context.Context, orderID, userID, driverID, status string) error {
	order, err := s.orderStorage.Get(orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return ErrNotYourOrder
	}

	if driverID != "" {
		if err = s.tripSaver.UpdateDriver(ctx, orderID, driverID); err != nil {
			return err
		}
	}

	return s.orderStorage.Update(orderID, status)
}

func (s *OrderService) Delete(orderID, userID string) error {
	order, err := s.orderStorage.Get(orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return ErrNotYourOrder
	}

	return s.orderStorage.Delete(orderID)
}

func (s *OrderService) SendOrder(ctx context.Context, userID string, order *models.Order) (string, error) {
	reqCtx, cancel := context.WithTimeout(ctx, reqTimeout)
	defer cancel()

	orderID, err := s.client.CreateOrder(reqCtx, userID, order)
	if err != nil {
		return "", err
	}

	if err = s.tripSaver.Save(ctx, orderID, userID, order); err != nil {
		return "", err
	}

	userOrder := &models.UserOrder{
		UserID:   userID,
		TaxiType: order.TaxiType,
		From:     order.From,
		To:       order.To,
		Status:   models.ProcessStatus,
	}

	if err = s.orderStorage.Save(orderID, userOrder); err != nil {
		return "", err
	}

	return orderID, nil
}

func (s *OrderService) SendComment(ctx context.Context, orderID, comment string) error {
	reqCtx, cancel := context.WithTimeout(ctx, reqTimeout)
	defer cancel()

	return s.client.SendComment(reqCtx, orderID, comment)
}

func (s *OrderService) Get(orderID, userID string) (*models.UserOrder, error) {
	order, err := s.orderStorage.Get(orderID)
	if err != nil {
		return nil, err
	}

	if order.UserID != userID {
		return nil, ErrNotYourOrder
	}

	return order, nil
}
