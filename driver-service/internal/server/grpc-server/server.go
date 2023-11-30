package grpc_server

import (
	"fmt"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/config"
	driver_service "github.com/HeadGardener/protos/gen/driver"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct{}

func (s *GRPCServer) Init(conf config.GRPCServerConfig, handler driver_service.DriverServiceServer) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", conf.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %e", err)
	}

	grpcServer := grpc.NewServer()
	driver_service.RegisterDriverServiceServer(grpcServer, handler)

	return grpcServer.Serve(lis)
}
