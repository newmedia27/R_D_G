package chat

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"project/internal/models"
	"project/internal/repositories/common"
)

type lastMessageDocument struct {
	Text      string             `bson:"text"`
	UserID    primitive.ObjectID `bson:"user_id"`
	CreatedAt time.Time          `bson:"created_at"`
}

type document struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty"`
	Type        models.ChatType      `bson:"type"`
	Name        string               `bson:"name"`
	Description string               `bson:"description"`
	OwnerID     primitive.ObjectID   `bson:"owner_id"`
	Members     []primitive.ObjectID `bson:"members"`
	LastMessage *lastMessageDocument `bson:"last_message,omitempty"`
	CreatedAt   time.Time            `bson:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at"`
}

func toModel(d *document) *models.Chat {
	members := make([]string, len(d.Members))
	for i, member := range d.Members {
		members[i] = member.Hex()
	}
	return &models.Chat{
		ID:          d.ID.Hex(),
		Type:        d.Type,
		Name:        d.Name,
		Description: d.Description,
		OwnerID:     d.OwnerID.Hex(),
		Members:     members,
		LastMessage: toLastMessage(d.LastMessage),
		CreatedAt:   d.CreatedAt,
		UpdatedAt:   d.UpdatedAt,
	}
}

func toDocument(c *models.Chat) (*document, error) {
	id, err := common.ParseOrGeneratePrimitiveId(c.ID)
	if err != nil {
		return nil, err
	}
	ownerID, err := common.ParseOrGeneratePrimitiveId(c.OwnerID)
	if err != nil {
		return nil, err
	}
	members := make([]primitive.ObjectID, len(c.Members))
	for i, member := range c.Members {
		members[i], err = common.ParseOrGeneratePrimitiveId(member)
		if err != nil {
			return nil, err
		}
	}
	return &document{
		ID:          id,
		Type:        c.Type,
		Name:        c.Name,
		Description: c.Description,
		OwnerID:     ownerID,
		Members:     members,
		LastMessage: toLastMessageDocument(c.LastMessage),
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}, nil
}
func toLastMessageDocument(d *models.LastMessage) *lastMessageDocument {
	if d == nil {
		return nil
	}
	id, err := common.ParseOrGeneratePrimitiveId(d.UserID)
	if err != nil {
		return nil
	}
	return &lastMessageDocument{
		Text:      d.Text,
		UserID:    id,
		CreatedAt: d.CreatedAt,
	}
}

func toLastMessage(d *lastMessageDocument) *models.LastMessage {
	if d == nil {
		return nil
	}
	return &models.LastMessage{
		Text:      d.Text,
		UserID:    d.UserID.Hex(),
		CreatedAt: d.CreatedAt,
	}
}
