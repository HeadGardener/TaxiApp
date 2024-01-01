package grpc_client

import (
	"context"
	"log"

	user_service "github.com/HeadGardener/protos/gen/user"
	"google.golang.org/grpc"
)

type UserServiceClient struct {
	c user_service.UserServiceClient
}

func NewUserServiceClient(conn *grpc.ClientConn) *UserServiceClient {
	c := user_service.NewUserServiceClient(conn)

	return &UserServiceClient{c: c}
}

func (c *UserServiceClient) AcceptOrder(ctx context.Context, userID, orderID, driverID string, status int32) error {
	req := &user_service.AcceptOrderRequest{
		UserID:   userID,
		OrderID:  orderID,
		DriverID: driverID,
		Status:   user_service.AcceptStatus(status),
	}

	resp, err := c.c.AcceptOrder(ctx, req)
	if err != nil {
		return err
	}

	log.Printf("user (%s) accepted order \n", resp.UserID)

	return nil
}

func (c *UserServiceClient) CompleteOrder(ctx context.Context, userID, orderID string, status int32) error {
	req := &user_service.CompleteOrderRequest{
		UserID:  userID,
		OrderID: orderID,
		Status:  user_service.CompleteStatus(status),
	}

	resp, err := c.c.CompleteOrder(ctx, req)
	if err != nil {
		return err
	}

	log.Printf("user (%s) completed order \n", resp.UserID)

	return nil
}
