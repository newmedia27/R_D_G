package chat

import (
	"time"

	"project/internal/models"
)

type CreateGroupRequest struct {
	Name        string   `json:"name"        validate:"required,min=3,max=50"`
	Description string   `json:"description" validate:"max=200"`
	Members     []string `json:"members"`
}

type ChatResponse struct {
	ID          string              `json:"id"`
	Type        models.ChatType     `json:"type"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	OwnerID     string              `json:"owner_id"`
	Members     []string            `json:"members"`
	LastMessage *models.LastMessage `json:"last_message"`
	CreatedAt   time.Time           `json:"created_at"`
	UpdatedAt   time.Time           `json:"updated_at"`
}

type AddMemberRequest struct {
	UserID string `json:"user_id" validate:"required"`
}
