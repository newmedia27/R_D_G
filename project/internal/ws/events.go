package ws

import "project/internal/models"

type EventType string

const (
	//FE->BE
	EventMessageSend    EventType = "message.send"
	EventMessageEdit    EventType = "message.edit"
	EventMessageDelete  EventType = "message.delete"
	EventUserTyping     EventType = "user.typing"
	EventUserStopTyping EventType = "user.stop_typing"

	//BE->FE
	EventMessageNew     EventType = "message.new"
	EventMessageEdited  EventType = "message.edited"
	EventMessageDeleted EventType = "message.deleted"
	EventTyping         EventType = "typing"
	EventStopTyping     EventType = "stop_typing"
	EventError          EventType = "error"
	EventChatCreated    EventType = "chat.created"

	EventAuthExpiringSoon EventType = "auth.expiring_soon"

	EventCallInvite EventType = "call.invite"
	EventCallAccept EventType = "call.accept"
	EventCallReject EventType = "call.reject"
	EventCallEnd    EventType = "call.end"
)

// FE->
type IncomingEvent struct {
	Type      EventType `json:"type"`
	ChatID    string    `json:"chat_id,omitempty"`
	MessageID string    `json:"message_id,omitempty"`
	Text      string    `json:"text,omitempty"`

	RoomID string `json:"room_id,omitempty"`
}

// BE->
type OutgoingEvent struct {
	Type      EventType    `json:"type"`
	ChatID    string       `json:"chat_id,omitempty"`
	MessageID string       `json:"message_id,omitempty"`
	UserID    string       `json:"user_id,omitempty"`
	Message   any          `json:"message,omitempty"`
	Error     string       `json:"error,omitempty"`
	Chat      *models.Chat `json:"chat,omitempty"`

	RoomID string `json:"room_id,omitempty"`
}
