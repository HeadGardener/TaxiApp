package grpc_server

import (
	"fmt"
	"net"

	"github.com/HeadGardener/TaxiApp/order-service/internal/config"
	order_service "github.com/HeadGardener/protos/gen/order"
	"google.golang.org/grpc"
)

type GRPCServer struct{}

func (s *GRPCServer) Init(conf config.GRPCServerConfig, handler order_service.OrderServiceServer) error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", conf.Port))
	if err != nil {
		return fmt.Errorf("failed to listen: %e", err)
	}

	grpcServer := grpc.NewServer()
	order_service.RegisterOrderServiceServer(grpcServer, handler)

	return grpcServer.Serve(lis)
}
