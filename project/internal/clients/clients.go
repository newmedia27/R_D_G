package clients

import (
	"context"
	"fmt"

	"project/internal/clients/mongodb"
	"project/internal/config"
)

type Clients struct {
	Mongo *mongodb.Client
}

var (
	ErrConnectToMongo = "Error connect to mongo"
)

func NewClients(ctx context.Context, cfg *config.Config) (*Clients, error) {
	client, err := mongodb.NewClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("%s", ErrConnectToMongo)
	}

	return &Clients{
		Mongo: client,
	}, nil
}
