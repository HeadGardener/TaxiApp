package grpc_handler

import (
	"context"
	"github.com/HeadGardener/TaxiApp/order-service/internal/models"
	order_service "github.com/HeadGardener/protos/gen/order"
)

func (h *ProcessOrderHandler) CreateOrder(ctx context.Context,
	req *order_service.CreateOrderRequest) (*order_service.CreateOrderResponse, error) {
	order := models.Order{
		UserID:   req.UserID,
		From:     req.From,
		To:       req.To,
		TaxiType: models.TaxiType(req.TaxiType),
	}

	orderID, err := h.orderService.CreateOrder(ctx, order)
	if err != nil {
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
