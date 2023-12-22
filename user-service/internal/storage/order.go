package storage

import (
	"errors"
	"sync"

	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
)

const (
	orderStorageStartSize = 100
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
		orders: make(map[string]*models.UserOrder, orderStorageStartSize),
	}
}

func (s *OrderStorage) Save(orderID string, order *models.UserOrder) error {
	s.mu.Unlock()
	defer s.mu.Unlock()

	s.orders[orderID] = order

	return nil
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
