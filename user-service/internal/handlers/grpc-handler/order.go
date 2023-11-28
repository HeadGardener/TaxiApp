package grpc_handler

import (
	"context"
	user_service "github.com/HeadGardener/protos/gen"
)

func (h *ProcessOrderHandler) AcceptOrder(_ context.Context,
	req *user_service.AcceptOrderRequest) (*user_service.AcceptOrderResponse, error) {
	if err := h.orderService.Update(req.OrderID, req.UserID, req.Status.String()); err != nil {
		return &user_service.AcceptOrderResponse{
			UserID:    req.UserID,
			Confirmed: false,
		}, err
	}

	return &user_service.AcceptOrderResponse{
		UserID:    req.UserID,
		Confirmed: true,
	}, nil
}

func (h *ProcessOrderHandler) CompleteOrder(_ context.Context,
	req *user_service.CompleteOrderRequest) (*user_service.CompleteOrderResponse, error) {
	if req.Status == user_service.CompleteStatus_CANCELED {
		if err := h.orderService.Update(req.OrderID, req.UserID, req.Status.String()); err != nil {
			return &user_service.CompleteOrderResponse{
				UserID:    req.UserID,
				Confirmed: false,
			}, err
		}

		return &user_service.CompleteOrderResponse{
			UserID:    req.UserID,
			Confirmed: true,
		}, nil
	}

	if err := h.orderService.Delete(req.OrderID, req.UserID); err != nil {
		return &user_service.CompleteOrderResponse{
			UserID:    req.UserID,
			Confirmed: false,
		}, err
	}

	return &user_service.CompleteOrderResponse{
		UserID:    req.UserID,
		Confirmed: true,
	}, nil
}
