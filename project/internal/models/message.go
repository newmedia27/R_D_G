package models

import "time"

type MessageType string

const (
	MessageTypeText MessageType = "text"
	MessageTypeFile MessageType = "file"
)

type Message struct {
	ID        string      `json:"id"`
	ChatID    string      `json:"chat_id"`
	UserID    string      `json:"user_id"`
	Type      MessageType `json:"type"`
	Text      string      `json:"text"`
	FileID    string      `json:"file_id"`
	IsEdited  bool        `json:"is_edited"`
	IsDeleted bool        `json:"is_deleted"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
}
