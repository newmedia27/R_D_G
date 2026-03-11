package ws

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/gofiber/contrib/v3/websocket"
	"project/internal/logger"
	"project/internal/models"
	"project/internal/services"
)

type Handler struct {
	hub      *Hub
	services *services.Services
}

func NewHandler(services *services.Services, hub *Hub) *Handler {
	return &Handler{
		hub:      hub,
		services: services,
	}
}

func (h *Handler) Handle(c *websocket.Conn) {
	UserID, ok := c.Locals("user_id").(string)
	if !ok {
		return
	}
	tokenExp, ok := c.Locals("token_exp").(time.Time)
	if !ok || tokenExp.IsZero() {
		return
	}

	client := NewClient(c, UserID, tokenExp)
	h.hub.Register(client)

	go client.writePump()

	client.readPump(h.hub, h.onMessage)
}

func (h *Handler) onMessage(c *Client, data []byte) {
	var event IncomingEvent
	if err := json.Unmarshal(data, &event); err != nil {
		h.sendError(c, "invalid event format")
		return
	}

	switch event.Type {
	case EventMessageSend:
		h.handleMessageSend(c, event)
	case EventMessageEdit:
		h.handleMessageEdit(c, event)
	case EventMessageDelete:
		h.handleMessageDelete(c, event)
	case EventUserTyping:
		h.handleTyping(c, event)
	case EventUserStopTyping:
		h.handleStopTyping(c, event)
	default:
		h.sendError(c, "unknown event type")
	}
}

func (h *Handler) sendError(c *Client, msg string) {
	data, err := json.Marshal(OutgoingEvent{
		Type:  EventError,
		Error: msg,
	})
	if err != nil {
		logger.LogErrorContext(context.Background(), err)
		return
	}
	select {
	case c.send <- data:
	default:
		//
	}
}

func (h *Handler) handleMessageSend(c *Client, e IncomingEvent) {
	if e.ChatID == "" {
		h.sendError(c, "chat id is required")
		return
	}
	if e.Text == "" {
		h.sendError(c, "text is required")
		return
	}
	ctx := context.Background()
	msg, err := h.services.Message.Create(ctx, e.ChatID, c.UserID, e.Text)
	if err != nil {
		if errors.Is(err, models.ErrForbidden) {
			h.sendError(c, "forbidden")
			return
		}
		h.sendError(c, "internal server error")
		logger.LogErrorContext(context.Background(), err)
		return
	}

	chat, err := h.services.Chat.GetChat(ctx, e.ChatID, c.UserID)
	if err != nil {
		h.sendError(c, "internal server error")
		logger.LogErrorContext(ctx, err)
		return
	}
	data, err := json.Marshal(OutgoingEvent{
		Type:    EventMessageNew,
		ChatID:  e.ChatID,
		Message: msg,
	})

	if err != nil {
		logger.LogErrorContext(ctx, err)
		return
	}
	h.hub.Broadcast(&BroadcastMsg{
		chat.Members, data,
	})
}

func (h *Handler) handleMessageEdit(c *Client, e IncomingEvent) {
	if e.MessageID == "" {
		h.sendError(c, "message_id is required")
		return
	}
	if e.Text == "" {
		h.sendError(c, "text is required")
		return
	}
	ctx := context.Background()
	msg, err := h.services.Message.EditMessage(ctx, e.MessageID, c.UserID, e.Text)
	if err != nil {
		if errors.Is(err, models.ErrForbidden) {
			h.sendError(c, "forbidden")
			return
		}
		if errors.Is(err, models.ErrMessageNotFound) {
			h.sendError(c, "message not found")
			return
		}
		h.sendError(c, "internal server error")
		logger.LogErrorContext(context.Background(), err)
		return
	}
	chat, err := h.services.Chat.GetChat(context.Background(), msg.ChatID, c.UserID)
	if err != nil {
		h.sendError(c, "internal server error")
		logger.LogErrorContext(context.Background(), err)
		return
	}
	data, err := json.Marshal(OutgoingEvent{
		Type:    EventMessageEdited,
		ChatID:  msg.ChatID,
		Message: msg,
	})
	if err != nil {
		logger.LogErrorContext(context.Background(), err)
		return
	}
	h.hub.Broadcast(&BroadcastMsg{
		chat.Members, data,
	})
}

func (h *Handler) handleMessageDelete(c *Client, e IncomingEvent) {
	if e.MessageID == "" {
		h.sendError(c, "message_id is required")
		return
	}

	ctx := context.Background()
	msg, err := h.services.Message.GetMessage(ctx, e.MessageID)
	if err != nil {
		if errors.Is(err, models.ErrMessageNotFound) {
			h.sendError(c, "message not found")
			return
		}
		h.sendError(c, "internal server error")
		logger.LogErrorContext(context.Background(), err)
		return
	}
	chat, err := h.services.Chat.GetChat(context.Background(), msg.ChatID, c.UserID)
	if err != nil {
		h.sendError(c, "internal server error")
		logger.LogErrorContext(context.Background(), err)
		return
	}
	if err = h.services.Message.DeleteMessage(context.Background(), e.MessageID, c.UserID); err != nil {
		if errors.Is(err, models.ErrForbidden) {
			h.sendError(c, "forbidden")
			return
		}
		h.sendError(c, "internal server error")
		logger.LogErrorContext(context.Background(), err)
		return
	}
	data, err := json.Marshal(OutgoingEvent{
		Type:      EventMessageDeleted,
		ChatID:    msg.ChatID,
		MessageID: e.MessageID,
	})
	if err != nil {
		logger.LogErrorContext(context.Background(), err)
		return
	}
	h.hub.Broadcast(&BroadcastMsg{
		chat.Members, data,
	})
}

func (h *Handler) handleTyping(c *Client, event IncomingEvent) {
	if event.ChatID == "" {
		h.sendError(c, "chat_id is required")
		return
	}

	chat, err := h.services.Chat.GetChat(context.Background(), event.ChatID, c.UserID)
	if err != nil {
		if errors.Is(err, models.ErrForbidden) {
			h.sendError(c, "forbidden")
			return
		}
		h.sendError(c, "internal server error")
		return
	}

	data, err := json.Marshal(OutgoingEvent{
		Type:   EventTyping,
		ChatID: event.ChatID,
		UserID: c.UserID,
	})
	if err != nil {
		return
	}

	members := make([]string, 0, len(chat.Members)-1)
	for _, m := range chat.Members {
		if m != c.UserID {
			members = append(members, m)
		}
	}
	h.hub.Broadcast(&BroadcastMsg{
		chat.Members, data,
	})
}

func (h *Handler) handleStopTyping(c *Client, event IncomingEvent) {
	if event.ChatID == "" {
		h.sendError(c, "chat_id is required")
		return
	}

	chat, err := h.services.Chat.GetChat(context.Background(), event.ChatID, c.UserID)
	if err != nil {
		if errors.Is(err, models.ErrForbidden) {
			h.sendError(c, "forbidden")
			return
		}
		h.sendError(c, "internal server error")
		return
	}

	data, err := json.Marshal(OutgoingEvent{
		Type:   EventStopTyping,
		ChatID: event.ChatID,
		UserID: c.UserID,
	})
	if err != nil {
		return
	}

	members := make([]string, 0, len(chat.Members)-1)
	for _, m := range chat.Members {
		if m != c.UserID {
			members = append(members, m)
		}
	}
	h.hub.Broadcast(&BroadcastMsg{
		chat.Members, data,
	})
}
