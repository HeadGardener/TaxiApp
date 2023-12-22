package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/HeadGardener/TaxiApp/user-service/internal/lib/auth"

	"github.com/redis/go-redis/v9"
)

type TokenStorage struct {
	rdb *redis.Client
}

func NewTokenStorage(rdb *redis.Client) *TokenStorage {
	return &TokenStorage{rdb: rdb}
}

func (s *TokenStorage) Add(ctx context.Context, userID, token string) error {
	s.rdb.Del(ctx, userID)
	err := s.rdb.Set(ctx, userID, token, auth.TokenTTL).Err()
	if err != nil {
		return fmt.Errorf("unable to store token: %w", err)
	}

	return nil
}

func (s *TokenStorage) Check(ctx context.Context, userID, token string) error {
	t, err := s.rdb.Get(ctx, userID).Result()
	if err != nil {
		return fmt.Errorf("user session doesn't exist: %w", err)
	}

	if t != token {
		return errors.New("tokens are different")
	}

	return nil
}

func (s *TokenStorage) Delete(ctx context.Context, userID string) error {
	if err := s.rdb.Del(ctx, userID).Err(); err != nil {
		return err
	}

	return nil
}
