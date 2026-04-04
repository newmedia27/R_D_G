package message

import (
	"time"

	"project/internal/models"
)

type SendMessageRequest struct {
	Text string `json:"text" validate:"required,min=1,max=4096"`
}

type EditMessageRequest struct {
	Text string `json:"text" validate:"required,min=1,max=4096"`
}

type MessageResponse struct {
	ID        string             `json:"id"`
	ChatID    string             `json:"chat_id"`
	UserID    string             `json:"user_id"`
	Type      models.MessageType `json:"type"`
	Text      string             `json:"text"`
	IsEdited  bool               `json:"is_edited"`
	IsDeleted bool               `json:"is_deleted"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}
