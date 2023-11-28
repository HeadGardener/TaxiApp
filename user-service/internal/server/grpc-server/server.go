package grpc_server

import (
	"fmt"
	"github.com/HeadGardener/TaxiApp/user-service/internal/config"
	user_service "github.com/HeadGardener/protos/gen"
	"google.golang.org/grpc"
	"net"
)

type GRPCServer struct{}

func (s *GRPCServer) Init(conf config.GRPCServerConfig, handler user_service.UserServiceServer) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", conf.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %e", err)
	}

	grpcServer := grpc.NewServer()
	user_service.RegisterUserServiceServer(grpcServer, handler)
	err = grpcServer.Serve(lis)
	if err != nil {
		return fmt.Errorf("failed to start grpc server: %e", err)
	}

	return nil
}
