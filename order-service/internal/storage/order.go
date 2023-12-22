package storage

import (
	"context"

	"github.com/HeadGardener/TaxiApp/order-service/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderStorage struct {
	coll *mongo.Collection
}

func (s *OrderStorage) Save(ctx context.Context, order *models.Order) (string, error) {
	_, err := s.coll.InsertOne(ctx, order)
	if err != nil {
		return "", err
	}

	return order.ID, err
}

func (s *OrderStorage) GetByID(ctx context.Context, orderID string) (models.Order, error) {
	var order models.Order
	if err := s.coll.FindOne(ctx, bson.D{{"id", orderID}}).Decode(&order); err != nil {
		return models.Order{}, err
	}

	return order, nil
}

func (s *OrderStorage) AddComment(ctx context.Context, orderID, comment string) error {
	filter := bson.D{{"id", orderID}}
	update := bson.D{{"$set", bson.D{{"comment", comment}}}}

	if _, err := s.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (s *OrderStorage) UpdateStatus(ctx context.Context, orderID, status string) error {
	filter := bson.D{{"id", orderID}}
	update := bson.D{{"$set", bson.D{{"status", status}}}}

	if _, err := s.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}

func (s *OrderStorage) UpdateRating(ctx context.Context, orderID string, rating float64) error {
	filter := bson.D{{"id", orderID}}
	update := bson.D{{"$set", bson.D{{"rating", rating}}}}

	if _, err := s.coll.UpdateOne(ctx, filter, update); err != nil {
		return err
	}

	return nil
}
