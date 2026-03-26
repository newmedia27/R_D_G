package livekit

import (
	"github.com/livekit/protocol/auth"
	"project/internal/config"
)

type Client struct {
	apiKey    string
	apiSecret string
	url       string
}

func NewClient(cfg *config.Config) *Client {
	return &Client{
		apiKey:    cfg.LiveKitAPIKey,
		apiSecret: cfg.LiveKitAPISecret,
		url:       cfg.LiveKitURL,
	}
}

func (c *Client) GenerateToken(userID, roomName string) (string, error) {
	authToken := auth.NewAccessToken(c.apiKey, c.apiSecret)

	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     roomName,
	}
	authToken.AddGrant(grant).SetIdentity(userID)
	return authToken.ToJWT()
}
