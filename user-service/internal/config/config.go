package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig         DBConfig
	ServerConfig     ServerConfig
	GRPCServerConfig GRPCServerConfig
	GRPCClientConfig GRPCClientConfig
	RedisConfig      RedisConfig
}

type DBConfig struct {
	URL string
}

type ServerConfig struct {
	Port string
}

type GRPCServerConfig struct {
	Port string
}

type GRPCClientConfig struct {
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

	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		return nil, errors.New("db url is empty")
	}

	srvport := os.Getenv("SERVER_PORT")
	if srvport == "" {
		return nil, errors.New("server port is empty")
	}

	grpcport := os.Getenv("GRPC_SERVER_PORT")
	if grpcport == "" {
		return nil, errors.New("grpc server port is empty")
	}

	grpcclientport := os.Getenv("GRPC_CLIENT_PORT")
	if grpcclientport == "" {
		return nil, errors.New("grpc client port is empty")
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
		GRPCServerConfig: GRPCServerConfig{
			Port: grpcport,
		},
		GRPCClientConfig: GRPCClientConfig{
			Port: grpcclientport,
		},
		RedisConfig: RedisConfig{
			Addr:     redisAddr,
			Password: redisPassword,
			DB:       redisDB,
		},
	}, nil
}
