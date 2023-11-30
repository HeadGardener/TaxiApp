package grpc_handler

import (
	"context"
	user_service "github.com/HeadGardener/protos/gen/user"
)

type OrderService interface {
	Update(ctx context.Context, orderID, userID, driverID, status string) error
	Delete(orderID, userID string) error
}

type ProcessOrderHandler struct {
	user_service.UnimplementedUserServiceServer

	orderService OrderService
}

func NewOrderHandler(orderService OrderService) *ProcessOrderHandler {
	return &ProcessOrderHandler{
		orderService: orderService,
	}
}
