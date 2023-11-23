package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerConfig ServerConfig

	DBConfig    DBConfig
	RedisConfig RedisConfig
}

type DBConfig struct {
	URL string
}

type ServerConfig struct {
	Port string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

func Init(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		return nil, fmt.Errorf("invalid path: %s", path)
	}

	srvport := os.Getenv("SERVER_PORT")
	if srvport == "" {
		return nil, errors.New("server port is empty")
	}

	dburl := os.Getenv("DB_URL")
	if dburl == "" {
		return nil, errors.New("db name is empty")
	}

	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		return nil, errors.New("redis addr is empty")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, errors.New("redis db is empty")
	}

	return &Config{
		DBConfig: DBConfig{
			URL: dburl,
		},
		ServerConfig: ServerConfig{
			Port: srvport,
		},
		RedisConfig: RedisConfig{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       redisDB,
		},
	}, nil
}
