package config

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DBConfig         DBConfig
	ServerConfig     ServerConfig
	GRPCServerConfig GRPCServerConfig
	RedisConfig      RedisConfig
}

type DBConfig struct {
	DBName     string
	Collection string
	URL        string
}

type ServerConfig struct {
	Port string
}

type GRPCServerConfig struct {
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

	dbname := os.Getenv("DATABASE_NAME")
	if dbname == "" {
		return nil, errors.New("db name is empty")
	}

	coll := os.Getenv("COLLECTION")
	if coll == "" {
		return nil, errors.New("collection is empty")
	}

	dburl := os.Getenv("DATABASE_URL")
	if dburl == "" {
		return nil, errors.New("db name is empty")
	}

	srvport := os.Getenv("SERVER_PORT")
	if srvport == "" {
		return nil, errors.New("server port is empty")
	}

	grpcport := os.Getenv("GRPC_PORT")
	if dburl == "" {
		return nil, errors.New("grpc port name is empty")
	}

	return &Config{
		DBConfig: DBConfig{
			DBName:     dbname,
			Collection: coll,
			URL:        dburl,
		},
		ServerConfig: ServerConfig{
			Port: srvport,
		},
		GRPCServerConfig: GRPCServerConfig{
			Port: grpcport,
		},
	}, nil
}
