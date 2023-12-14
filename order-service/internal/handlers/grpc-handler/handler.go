package grpc_handler

import (
	"context"
	"github.com/HeadGardener/TaxiApp/order-service/internal/models"
	order_service "github.com/HeadGardener/protos/gen/order"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order models.Order) (string, error)
	AddComment(ctx context.Context, orderID, comment string) error
}

type ProcessOrderHandler struct {
	order_service.UnimplementedOrderServiceServer

	orderService OrderService
}

func NewProcessOrderHandler(orderService OrderService) *ProcessOrderHandler {
	return &ProcessOrderHandler{orderService: orderService}
}
