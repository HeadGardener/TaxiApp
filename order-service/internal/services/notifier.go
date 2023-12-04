package services

import (
	"context"
	"errors"
	"github.com/HeadGardener/TaxiApp/order-service/internal/models"
	"log"
	"time"
)

var (
	ErrAddTimeout      = errors.New("end of adding timeout")
	ErrInvalidTaxiType = errors.New("invalid taxi type")
)

var (
	waitTime = 5 * time.Minute
)

type DriverServiceGRPCClient interface {
	ConsumeOrder(ctx context.Context, driverID string, order models.OrderInfo) error
}

type OrderNotifier struct {
	client        DriverServiceGRPCClient
	usersQueues   *models.UsersQueues
	driversQueues *models.DriversQueues
}

func NewOrderNotifier(client DriverServiceGRPCClient, usersQueues *models.UsersQueues, driversQueues *models.DriversQueues) *OrderNotifier {
	return &OrderNotifier{
		client:        client,
		usersQueues:   usersQueues,
		driversQueues: driversQueues,
	}
}

func (n *OrderNotifier) AddUserToQueue(ctx context.Context, userID, orderID, from, to string, taxiType models.TaxiType) error {
	var err error

	go func() {
		for {
			select {
			case <-ctx.Done():
				err = errors.Join(ErrAddTimeout, ctx.Err())
				return
			}
		}
	}()

	order := models.OrderInfo{
		UserID:  userID,
		OrderID: orderID,
		From:    from,
		To:      to,
	}

	switch taxiType {
	case models.Economy:
		n.usersQueues.Economy <- order

	case models.Comfort:
		n.usersQueues.Comfort <- order

	case models.Business:
		n.usersQueues.Business <- order

	default:
		err = ErrInvalidTaxiType
	}

	return err
}

func (n *OrderNotifier) AddDriverToQueue(ctx context.Context, driverID string, taxiType models.TaxiType) error {
	var err error

	go func() {
		for {
			select {
			case <-ctx.Done():
				err = errors.Join(ErrAddTimeout, ctx.Err())
				return
			}
		}
	}()

	switch taxiType {
	case models.Economy:
		n.driversQueues.Economy <- driverID

	case models.Comfort:
		n.driversQueues.Comfort <- driverID

	case models.Business:
		n.driversQueues.Business <- driverID

	default:
		err = ErrInvalidTaxiType
	}

	return err
}

func (n *OrderNotifier) MakeUpOrders() {
	done := make(chan struct{}, 1)
	go func() {
	usersLoop:
		for {
			order := <-n.usersQueues.Economy
			log.Printf("[INFO] finding order for user %s", order.UserID)
			ctx, cancel := context.WithTimeout(context.Background(), waitTime)

		driversLoop:
			for {
				select {
				case driverID := <-n.driversQueues.Economy:
					err := n.client.ConsumeOrder(ctx, driverID, order)
					if err != nil {
						log.Printf("[ERROR] failed while sending driver accept order request: %s", err.Error())
					}
					continue driversLoop

				case <-ctx.Done():
					// send request for user that his time of waiting is over
					cancel()
					continue usersLoop
				}
			}
		}
	}()

	go func() {
	usersLoop:
		for {
			order := <-n.usersQueues.Comfort
			log.Printf("[INFO] finding order for user %s", order.UserID)
			ctx, cancel := context.WithTimeout(context.Background(), waitTime)

		driversLoop:
			for {
				select {
				case driverID := <-n.driversQueues.Comfort:
					err := n.client.ConsumeOrder(ctx, driverID, order)
					if err != nil {
						log.Printf("[ERROR] failed while sending driver accept order request: %s", err.Error())
					}
					continue driversLoop

				case <-ctx.Done():
					// send request for user that his time of waiting is over
					cancel()
					continue usersLoop
				}
			}
		}
	}()

	go func() {
	usersLoop:
		for {
			order := <-n.usersQueues.Business
			log.Printf("[INFO] finding order for user %s", order.UserID)
			ctx, cancel := context.WithTimeout(context.Background(), waitTime)

		driversLoop:
			for {
				select {
				case driverID := <-n.driversQueues.Business:
					err := n.client.ConsumeOrder(ctx, driverID, order)
					if err != nil {
						log.Printf("[ERROR] failed while sending driver accept order request: %s", err.Error())
					}
					continue driversLoop

				case <-ctx.Done():
					// send request for user that his time of waiting is over
					cancel()
					continue usersLoop
				}
			}
		}
	}()

	<-done
}
