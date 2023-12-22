package grpc_handler

import (
	"context"

	"github.com/HeadGardener/TaxiApp/order-service/internal/models"
	order_service "github.com/HeadGardener/protos/gen/order"
)

func (h *ProcessOrderHandler) CreateOrder(ctx context.Context,
	req *order_service.CreateOrderRequest) (*order_service.CreateOrderResponse, error) {
	order := &models.Order{
		UserID:   req.UserID,
		From:     req.From,
		To:       req.To,
		TaxiType: models.TaxiType(req.TaxiType),
	}

	orderID, err := h.orderService.CreateOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	if err = h.orderNotifier.AddUserToQueue(ctx,
		req.UserID,
		orderID,
		req.From,
		req.To,
		models.TaxiType(req.TaxiType)); err != nil {
		return nil, err
	}

	return &order_service.CreateOrderResponse{
		OrderID: orderID,
	}, nil
}

func (h *ProcessOrderHandler) AddComment(ctx context.Context,
	req *order_service.AddCommentRequest) (*order_service.AddCommentResponse, error) {
	if err := h.orderService.AddComment(ctx, req.OrderID, req.Comment); err != nil {
		return nil, err
	}

	return &order_service.AddCommentResponse{
		Added: true,
	}, nil
}

func (h *ProcessOrderHandler) AddDriver(ctx context.Context,
	req *order_service.AddDriverRequest) (*order_service.AddCommentResponse, error) {
	if err := h.orderNotifier.AddDriverToQueue(ctx, req.DriverID, models.TaxiType(req.TaxiType)); err != nil {
		return nil, err
	}

	return &order_service.AddCommentResponse{
		Added: true,
	}, nil
}

func (h *ProcessOrderHandler) ProcessOrder(ctx context.Context,
	req *order_service.ProcessOrderRequest) (*order_service.ProcessOrderResponse, error) {
	if err := h.orderService.ProcessOrder(ctx, req.OrderID, req.DriverID, req.Status.String()); err != nil {
		return nil, err
	}

	return &order_service.ProcessOrderResponse{
		Processed: true,
	}, nil
}

func (h *ProcessOrderHandler) CompleteOrder(ctx context.Context,
	req *order_service.CompleteOrderRequest) (*order_service.CompleteOrderResponse, error) {
	if err := h.orderService.CompleteOrder(ctx,
		req.OrderID,
		req.DriverID,
		req.Status.String(),
		float64(req.Rating)); err != nil {
		return nil, err
	}

	return &order_service.CompleteOrderResponse{
		Completed: true,
	}, nil
}
