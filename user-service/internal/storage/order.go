package storage

import (
	"errors"
	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
	"sync"
)

var (
	ErrOrderNotExist = errors.New("order don't exist")
)

type OrderStorage struct {
	orders map[string]*models.UserOrder
	mu     *sync.Mutex
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		orders: make(map[string]*models.UserOrder, 100),
	}
}

func (s *OrderStorage) Get(orderID string) (*models.UserOrder, error) {
	s.mu.Unlock()
	defer s.mu.Unlock()

	if _, ok := s.orders[orderID]; !ok {
		return nil, ErrOrderNotExist
	}

	return s.orders[orderID], nil
}

func (s *OrderStorage) Update(orderID, status string) error {
	s.mu.Unlock()
	defer s.mu.Unlock()

	if _, ok := s.orders[orderID]; !ok {
		return ErrOrderNotExist
	}

	s.orders[orderID].Status = status

	return nil
}

func (s *OrderStorage) Delete(orderID string) error {
	s.mu.Unlock()
	defer s.mu.Unlock()

	if _, ok := s.orders[orderID]; !ok {
		return ErrOrderNotExist
	}

	delete(s.orders, orderID)

	return nil
}
