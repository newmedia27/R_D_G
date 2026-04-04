package message

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"project/internal/models"
	"project/internal/repositories/common"
)

type document struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ChatID    primitive.ObjectID `bson:"chat_id"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Type      models.MessageType `bson:"type"`
	Text      string             `bson:"text"`
	FileID    string             `bson:"file_id"`
	IsEdited  bool               `bson:"is_edited"`
	IsDeleted bool               `bson:"is_deleted"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

func toModel(d *document) *models.Message {
	return &models.Message{
		ID:        d.ID.Hex(),
		ChatID:    d.ChatID.Hex(),
		UserID:    d.UserID.Hex(),
		Type:      d.Type,
		Text:      d.Text,
		FileID:    d.FileID,
		IsEdited:  d.IsEdited,
		IsDeleted: d.IsDeleted,
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func toDocument(m *models.Message) (*document, error) {
	id, err := common.ParseOrGeneratePrimitiveId(m.ID)
	if err != nil {
		return nil, err
	}
	chatID, err := common.ParseOrGeneratePrimitiveId(m.ChatID)
	if err != nil {
		return nil, err
	}
	userID, err := common.ParseOrGeneratePrimitiveId(m.UserID)
	if err != nil {
		return nil, err
	}

	return &document{
		ID:        id,
		ChatID:    chatID,
		UserID:    userID,
		Type:      m.Type,
		Text:      m.Text,
		FileID:    m.FileID,
		IsEdited:  m.IsEdited,
		IsDeleted: m.IsDeleted,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}, nil
}
