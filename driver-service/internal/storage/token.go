package storage

import (
	"context"
	"errors"
	"fmt"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/lib/auth"
	"github.com/redis/go-redis/v9"
)

type TokenStorage struct {
	rdb *redis.Client
}

func NewTokenStorage(rdb *redis.Client) *TokenStorage {
	return &TokenStorage{rdb: rdb}
}

func (s *TokenStorage) Add(ctx context.Context, driverID, token string) error {
	s.rdb.Del(ctx, driverID)
	err := s.rdb.Set(ctx, driverID, token, auth.TokenTTL).Err()
	if err != nil {
		return fmt.Errorf("unable to store token: %w", err)
	}

	return nil
}

func (s *TokenStorage) Check(ctx context.Context, driverID, token string) error {
	t, err := s.rdb.Get(ctx, driverID).Result()
	if err != nil {
		return fmt.Errorf("user session doesn't exist: %w", err)
	}

	if t != token {
		return errors.New("tokens are different")
	}

	return nil
}

func (s *TokenStorage) Delete(ctx context.Context, driverID string) error {
	if err := s.rdb.Del(ctx, driverID).Err(); err != nil {
		return err
	}

	return nil
}
