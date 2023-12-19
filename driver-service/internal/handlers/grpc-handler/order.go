package grpc_handler

import (
	"context"
	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
	driver_service "github.com/HeadGardener/protos/gen/driver"
)

func (h *ProcessOrderHandler) ConsumeOrder(_ context.Context,
	req *driver_service.ConsumeOrderRequest) (*driver_service.ConsumeOrderResponse, error) {
	order := models.Order{
		ID:       req.OrderID,
		DriverID: req.DriverID,
		UserID:   req.UserID,
		From:     req.From,
		To:       req.To,
		Status:   models.ConsumedStatus,
	}

	if err := h.orderService.Add(order); err != nil {
		return nil, err
	}

	return &driver_service.ConsumeOrderResponse{
		DriverID: req.DriverID,
		OrderID:  req.OrderID,
	}, nil
}
