package services

import (
	"errors"
	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
)

var (
	ErrNotYourOrder = errors.New("this not your order")
)

type OrderStorage interface {
	Get(orderID string) (*models.UserOrder, error)
	Update(orderID, status string) error
	Delete(orderID string) error
}

type OrderService struct {
	orderStorage OrderStorage
}

func NewOrderService(orderStorage OrderStorage) *OrderService {
	return &OrderService{
		orderStorage: orderStorage,
	}
}

func (s *OrderService) Update(orderID, userID, status string) error {
	order, err := s.orderStorage.Get(orderID)
	if err != nil {
		return err
	}

	if order.UserID != userID {
		return ErrNotYourOrder
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
