package storage

import (
	"github.com/HeadGardener/TaxiApp/user-service/internal/config"
	"github.com/redis/go-redis/v9"
)

func NewRedisDB(conf config.RedisConfig) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.DB,
	})
}
