package storage

import (
	"errors"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
	"sync"
)

var (
	ErrOrderNotExist = errors.New("order don't exist")
)

type OrderStorage struct {
	orders map[string]models.Order
	mu     *sync.Mutex
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]models.Order, 100),
	}
}

func (s *OrderStorage) Add(order models.Order) error {
	s.mu.Unlock()
	defer s.mu.Unlock()

	s.orders[order.ID] = order

	return nil
}

func (s *OrderStorage) GetByDriverID(driverID string) (models.Order, error) {
	s.mu.Unlock()
	defer s.mu.Unlock()

	if _, ok := s.orders[driverID]; !ok {
		return models.Order{}, ErrOrderNotExist
	}

	return s.orders[driverID], nil
}

func (s *OrderStorage) Delete(driverID string) error {
	s.mu.Unlock()
	defer s.mu.Unlock()

	if _, ok := s.orders[driverID]; !ok {
		return ErrOrderNotExist
	}

	delete(s.orders, driverID)

	return nil
}
