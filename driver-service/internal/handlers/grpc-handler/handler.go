package grpc_handler

import (
	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
	driver_service "github.com/HeadGardener/protos/gen/driver"
)

type OrderService interface {
	Add(order models.Order) error
}

type ProcessOrderHandler struct {
	driver_service.UnimplementedDriverServiceServer

	orderService OrderService
}

func NewProcessOrderHandler(orderService OrderService) *ProcessOrderHandler {
	return &ProcessOrderHandler{orderService: orderService}
}
