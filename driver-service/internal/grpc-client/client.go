package grpc_client

import (
	"context"

	"github.com/HeadGardener/TaxiApp/driver-service/internal/models"
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

func (c *OrderServiceClient) AddDriver(ctx context.Context, driverID string, taxiType models.TaxiType) error {
	req := &order_service.AddDriverRequest{
		DriverID: driverID,
		TaxiType: int32(taxiType),
	}

	_, err := c.c.AddDriver(ctx, req)

	return err
}

func (c *OrderServiceClient) ProcessOrder(ctx context.Context, driverID, orderID string,
	status models.AcceptOrderStatus) error {
	req := &order_service.ProcessOrderRequest{
		DriverID: driverID,
		OrderID:  orderID,
		Status:   order_service.AcceptStatus(status),
	}

	_, err := c.c.ProcessOrder(ctx, req)

	return err
}

func (c *OrderServiceClient) CompleteOrder(ctx context.Context, driverID, orderID string,
	status models.CompleteOrderStatus, rating float64) error {
	req := &order_service.CompleteOrderRequest{
		DriverID: driverID,
		OrderID:  orderID,
		Status:   order_service.CompleteStatus(status),
		Rating:   float32(rating),
	}

	_, err := c.c.CompleteOrder(ctx, req)

	return err
}
