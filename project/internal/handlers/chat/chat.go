package chat

import (
	"encoding/json"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"project/internal/logger"
	"project/internal/models"
	"project/internal/services"
	"project/internal/ws"
)

type Handler struct {
	services  *services.Services
	validator *validator.Validate
	hub       *ws.Hub
}

func NewHandler(s *services.Services, v *validator.Validate, hub *ws.Hub) *Handler {
	return &Handler{
		services:  s,
		validator: v,
		hub:       hub,
	}
}

// CreateGroup
// @Summary Створити груповий чат
// @Tags chats
// @Accept json
// @Produce json
// @Param request body CreateGroupRequest true "Create group"
// @Success 201 {object} ChatResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/chats/group [post]
func (h *Handler) CreateGroup(c fiber.Ctx) error {
	var req CreateGroupRequest

	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}
	ownerID := c.Locals("user_id").(string)

	chat, err := h.services.Chat.CreateGroupChat(c.Context(), ownerID, req.Name, req.Description, req.Members)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	data, err := json.Marshal(ws.OutgoingEvent{
		Type: ws.EventChatCreated,
		Chat: chat,
	})
	if err != nil {
		logger.LogErrorContext(c.Context(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	members := make([]string, 0, len(chat.Members)-1)
	for _, m := range chat.Members {
		if m != ownerID {
			members = append(members, m)
		}
	}
	h.hub.Broadcast(&ws.BroadcastMsg{
		UserIDs: members,
		Data:    data,
	})

	return c.Status(fiber.StatusCreated).JSON(chat)
}

// CreatePrivate
// @Summary Створити приватний чат
// @Tags chats
// @Produce json
// @Param userID path string true "User ID"
// @Success 200 {object} ChatResponse
// @Failure 400 {object} map[string]interface{}
// @Router /api/v1/chats/private/{userID} [post]
func (h *Handler) CreatePrivate(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	targetID := c.Params("userID")

	if targetID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user id is required",
		})
	}

	if userID == targetID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "cannot create chat with yourself",
		})
	}

	chat, err := h.services.Chat.CreatePrivateChat(c.Context(), userID, targetID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	data, err := json.Marshal(ws.OutgoingEvent{
		Type: ws.EventChatCreated,
		Chat: chat,
	})
	if err != nil {
		logger.LogErrorContext(c.Context(), err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	members := make([]string, 0, len(chat.Members)-1)
	for _, m := range chat.Members {
		if m != userID {
			members = append(members, m)
		}
	}
	h.hub.Broadcast(&ws.BroadcastMsg{
		UserIDs: members,
		Data:    data,
	})

	return c.Status(fiber.StatusOK).JSON(chat)
}

// GetUserChats
// @Summary Отримати всі чати юзера
// @Tags chats
// @Produce json
// @Success 200 {array} ChatResponse
// @Router /api/v1/chats [get]
func (h *Handler) GetUserChats(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	chats, err := h.services.Chat.GetUserChats(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(chats)
}

// GetChat
// @Summary Отримати чат по id
// @Tags chats
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} ChatResponse
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/chats/{id} [get]
func (h *Handler) GetChat(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	chatID := c.Params("id")

	chat, err := h.services.Chat.GetChat(c.Context(), chatID, userID)
	if err != nil {
		if errors.Is(err, models.ErrChatNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "chat not found",
			})
		}
		if errors.Is(err, models.ErrForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "forbidden",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(chat)
}

// AddMember
// @Summary Додати учасника до групового чату
// @Tags chats
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param request body AddMemberRequest true "Add member"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/chats/{id}/members [post]
func (h *Handler) AddMember(c fiber.Ctx) error {
	ownerID := c.Locals("user_id").(string)
	chatID := c.Params("id")

	var req AddMemberRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.services.Chat.AddMember(c.Context(), chatID, ownerID, req.UserID); err != nil {
		if errors.Is(err, models.ErrChatNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "chat not found",
			})
		}
		if errors.Is(err, models.ErrForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "forbidden",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "member added",
	})
}

// RemoveMember
// @Summary Видалити учасника з групового чату
// @Tags chats
// @Produce json
// @Param id path string true "Chat ID"
// @Param userID path string true "User ID"
// @Success 200 {object} map[string]interface{}
// @Failure 403 {object} map[string]interface{}
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/chats/{id}/members/{userID} [delete]
func (h *Handler) RemoveMember(c fiber.Ctx) error {
	ownerID := c.Locals("user_id").(string)
	chatID := c.Params("id")
	userID := c.Params("userID")

	if err := h.services.Chat.RemoveMember(c.Context(), chatID, ownerID, userID); err != nil {
		if errors.Is(err, models.ErrChatNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "chat not found",
			})
		}
		if errors.Is(err, models.ErrForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "forbidden",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "member removed",
	})
}
