package ws

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

	EventAuthExpiringSoon EventType = "auth.expiring_soon"
)

// FE->
type IncomingEvent struct {
	Type      EventType `json:"type"`
	ChatID    string    `json:"chat_id,omitempty"`
	MessageID string    `json:"message_id,omitempty"`
	Text      string    `json:"text,omitempty"`
}

// BE->
type OutgoingEvent struct {
	Type      EventType `json:"type"`
	ChatID    string    `json:"chat_id,omitempty"`
	MessageID string    `json:"message_id,omitempty"`
	UserID    string    `json:"user_id,omitempty"`
	Message   any       `json:"message,omitempty"`
	Error     string    `json:"error,omitempty"`
}
