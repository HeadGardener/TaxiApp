package storage

import (
	"context"
	"github.com/HeadGardener/TaxiApp/order-service/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"time"
)

func NewMongoCollection(ctx context.Context, conf config.DBConfig) (*mongo.Collection, error) {
	ctxConn, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctxConn, options.Client().ApplyURI(conf.URL))
	if err != nil {
		return nil, err
	}

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = client.Ping(ctxPing, readpref.Primary())
	if err != nil {
		return nil, err
	}

	return client.Database(conf.DBName).Collection(conf.Collection), err
}
