package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBConfig               DBConfig
	ServerConfig           ServerConfig
	GRPCServerConfig       GRPCServerConfig
	GRPCUserClientConfig   GRPCUserClientConfig
	GRPCDriverClientConfig GRPCDriverClientConfig
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

type GRPCUserClientConfig struct {
	Port string
}

type GRPCDriverClientConfig struct {
	Port string
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

	grpcport := os.Getenv("GRPC_SERVER_PORT")
	if grpcport == "" {
		return nil, errors.New("grpc server port is empty")
	}

	grpcuserclientport := os.Getenv("GRPC_USER_CLIENT_PORT")
	if grpcuserclientport == "" {
		return nil, errors.New("grpc user client port is empty")
	}

	grpcdriverclientport := os.Getenv("GRPC_DRIVER_CLIENT_PORT")
	if grpcdriverclientport == "" {
		return nil, errors.New("grpc driver client port is empty")
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
		GRPCUserClientConfig: GRPCUserClientConfig{
			Port: grpcuserclientport,
		},
		GRPCDriverClientConfig: GRPCDriverClientConfig{
			Port: grpcdriverclientport,
		},
	}, nil
}
