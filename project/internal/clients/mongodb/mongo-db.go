package mongodb

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"project/internal/config"
)

type Client struct {
	Db *mongo.Database
}

func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	URI := fmt.Sprintf("mongodb://%s:%s@%s/%s?authSource=admin", cfg.MongoUser, cfg.MongoUserPassword, cfg.MongoURI, cfg.MongoDatabase)
	c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	cOptions := options.Client().ApplyURI(URI)
	client, err := mongo.Connect(ctx, cOptions)
	if err != nil {
		return nil, fmt.Errorf("can't connect to mongo: %w", err)
	}
	err = client.Ping(c, nil)
	if err != nil {
		return nil, fmt.Errorf("can't ping mongo: %w", err)
	}
	slog.Info("Connected to mongo")

	return &Client{
		Db: client.Database(cfg.MongoDatabase),
	}, nil
}

func DisconnectMongo(c *mongo.Database) {
	slog.Info("Disconnect")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if c != nil {
		if err := c.Client().Disconnect(ctx); err != nil {
			slog.Error("can't disconnect from mongo", slog.Any("err", err))
		}
	}
}
