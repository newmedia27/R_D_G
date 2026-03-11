package user

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"project/internal/models"
	"project/internal/repositories/common"
)

type document struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Username     string             `bson:"username,omitempty"`
	Email        string             `bson:"email,omitempty"`
	PasswordHash string             `bson:"password_hash,omitempty"`
	CreatedAt    time.Time          `bson:"created_at,omitempty"`
	UpdatedAt    time.Time          `bson:"updated_at,omitempty"`
}

func toModel(d *document) *models.User {
	return &models.User{
		ID:           d.ID.Hex(),
		Username:     d.Username,
		Email:        d.Email,
		PasswordHash: d.PasswordHash,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
	}
}

func toDocument(u *models.User) (*document, error) {
	id, err := common.ParseOrGeneratePrimitiveId(u.ID)
	if err != nil {
		return nil, err
	}
	return &document{
		ID:           id,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		CreatedAt:    u.CreatedAt,
		UpdatedAt:    u.UpdatedAt,
	}, nil

}
