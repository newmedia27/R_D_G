package session

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"project/internal/models"
	"project/internal/repositories/common"
)

type document struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	UserID           primitive.ObjectID `bson:"user_id"`
	RefreshTokenHash string             `bson:"refresh_token_hash"`
	DeviceSessionID  string             `bson:"device_session_id"`
	FingerPrint      string             `bson:"fingerprint"`
	UserAgent        string             `bson:"user_agent"`
	IP               string             `bson:"ip"`
	ExpiresAt        time.Time          `bson:"expires_at"`
	CreatedAt        time.Time          `bson:"created_at"`
	LastUsedAt       time.Time          `bson:"last_used_at"`
}

func toModel(doc *document) *models.Session {
	return &models.Session{
		ID:               doc.ID.Hex(),
		UserID:           doc.UserID.Hex(),
		RefreshTokenHash: doc.RefreshTokenHash,
		DeviceSessionID:  doc.DeviceSessionID,
		FingerPrint:      doc.FingerPrint,
		UserAgent:        doc.UserAgent,
		IP:               doc.IP,
		ExpiresAt:        doc.ExpiresAt,
		CreatedAt:        doc.CreatedAt,
		LastUsedAt:       doc.LastUsedAt,
	}
}

func toDocument(s *models.Session) (*document, error) {
	id, err := common.ParseOrGeneratePrimitiveId(s.ID)
	if err != nil {
		return nil, err
	}
	userID, err := primitive.ObjectIDFromHex(s.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	return &document{
		ID:               id,
		UserID:           userID,
		RefreshTokenHash: s.RefreshTokenHash,
		DeviceSessionID:  s.DeviceSessionID,
		FingerPrint:      s.FingerPrint,
		UserAgent:        s.UserAgent,
		IP:               s.IP,
		ExpiresAt:        s.ExpiresAt,
		CreatedAt:        s.CreatedAt,
		LastUsedAt:       s.LastUsedAt,
	}, nil
}
