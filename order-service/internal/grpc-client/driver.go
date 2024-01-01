package grpc_client

import (
	"context"
	"log"

	"github.com/HeadGardener/TaxiApp/order-service/internal/models"
	driver_service "github.com/HeadGardener/protos/gen/driver"
	"google.golang.org/grpc"
)

type DriverServiceClient struct {
	c driver_service.DriverServiceClient
}

func NewDriverServiceClient(conn *grpc.ClientConn) *DriverServiceClient {
	c := driver_service.NewDriverServiceClient(conn)

	return &DriverServiceClient{c: c}
}

func (c *DriverServiceClient) ConsumeOrder(ctx context.Context, driverID string, order models.OrderInfo) error {
	req := &driver_service.ConsumeOrderRequest{
		DriverID: driverID,
		UserID:   order.UserID,
		OrderID:  order.OrderID,
		From:     order.From,
		To:       order.To,
	}

	resp, err := c.c.ConsumeOrder(ctx, req)
	if err != nil {
		return err
	}

	log.Printf("order (%s) consumed by driver (%s) \n", resp.OrderID, resp.DriverID)

	return nil
}
