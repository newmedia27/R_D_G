package clients

import (
	"context"
	"fmt"

	"project/internal/clients/livekit"
	"project/internal/clients/mongodb"
	"project/internal/config"
)

type Clients struct {
	Mongo   *mongodb.Client
	LiveKit *livekit.Client
}

var (
	ErrConnectToMongo = "Error connect to mongo"
)

func NewClients(ctx context.Context, cfg *config.Config) (*Clients, error) {
	client, err := mongodb.NewClient(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("%s", ErrConnectToMongo)
	}

	liveKitClient := livekit.NewClient(cfg)

	return &Clients{
		Mongo:   client,
		LiveKit: liveKitClient,
	}, nil
}
