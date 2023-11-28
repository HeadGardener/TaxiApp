package grpc_handler

import user_service "github.com/HeadGardener/protos/gen"

type OrderService interface {
	Update(orderID, userID, status string) error
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
