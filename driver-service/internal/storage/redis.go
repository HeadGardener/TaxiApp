package storage

import (
	"github.com/HeadGardener/TaxiApp/driver-service/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisDB(config config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
}
