package grpc_handler

import (
	"context"
	"github.com/HeadGardener/TaxiApp/order-service/internal/models"
	order_service "github.com/HeadGardener/protos/gen/order"
)

type OrderService interface {
	CreateOrder(ctx context.Context, order models.Order) (string, error)
	AddComment(ctx context.Context, orderID, comment string) error

	ProcessOrder(ctx context.Context, orderID, driverID, status string) error
	CompleteOrder(ctx context.Context, driverID, orderID, status string, rating float64) error
}

type OrderNotifier interface {
	AddUserToQueue(ctx context.Context, userID, orderID, from, to string, taxiType models.TaxiType) error
	AddDriverToQueue(ctx context.Context, driverID string, taxiType models.TaxiType) error
}

type ProcessOrderHandler struct {
	order_service.UnimplementedOrderServiceServer

	orderService  OrderService
	orderNotifier OrderNotifier
}

func NewProcessOrderHandler(orderService OrderService, orderNotifier OrderNotifier) *ProcessOrderHandler {
	return &ProcessOrderHandler{
		orderService:  orderService,
		orderNotifier: orderNotifier,
	}
}
