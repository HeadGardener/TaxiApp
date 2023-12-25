package grpc_client

import (
	"context"

	"github.com/HeadGardener/TaxiApp/user-service/internal/models"
	order_service "github.com/HeadGardener/protos/gen/order"
	"google.golang.org/grpc"
)

type OrderServiceClient struct {
	c order_service.OrderServiceClient
}

func NewOrderServiceClient(conn *grpc.ClientConn) *OrderServiceClient {
	c := order_service.NewOrderServiceClient(conn)

	return &OrderServiceClient{c: c}
}

func (c *OrderServiceClient) CreateOrder(ctx context.Context, userID string, order *models.Order) (string, error) {
	req := &order_service.CreateOrderRequest{
		UserID:   userID,
		TaxiType: int32(order.TaxiType),
		From:     order.From,
		To:       order.To,
	}

	resp, err := c.c.CreateOrder(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.OrderID, nil
}

func (c *OrderServiceClient) SendComment(ctx context.Context, orderID, comment string) error {
	req := &order_service.AddCommentRequest{
		OrderID: orderID,
		Comment: comment,
	}

	_, err := c.c.AddComment(ctx, req)

	return err
}
