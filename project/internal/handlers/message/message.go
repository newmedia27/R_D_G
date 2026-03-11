package message

import (
	"errors"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"project/internal/models"
	"project/internal/services"
)

type Handler struct {
	services  *services.Services
	validator *validator.Validate
}

func NewHandler(s *services.Services, v *validator.Validate) *Handler {
	return &Handler{
		services:  s,
		validator: v,
	}
}

// GetMessages
// @Summary Отримати історію повідомлень
// @Tags messages
// @Produce json
// @Param id path string true "Chat ID"
// @Param before query string false "Message ID для пагінації"
// @Param limit query int false "Кількість повідомлень" default(50)
// @Success 200 {array} MessageResponse
// @Router /api/v1/chats/{id}/messages [get]
func (h *Handler) GetMessages(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	chatID := c.Params("id")
	before := c.Query("before", "")

	limitStr := c.Query("limit", "50")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 50
	}

	msgs, err := h.services.Message.GetMessages(c.Context(), chatID, userID, before, limit)
	if err != nil {
		if errors.Is(err, models.ErrForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "forbidden",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(msgs)
}

// EditMessage
// @Summary Редагувати повідомлення
// @Tags messages
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Param request body EditMessageRequest true "Edit message"
// @Success 200 {object} MessageResponse
// @Router /api/v1/messages/{id} [put]
func (h *Handler) EditMessage(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	messageID := c.Params("id")

	var req EditMessageRequest
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

	msg, err := h.services.Message.EditMessage(c.Context(), messageID, userID, req.Text)
	if err != nil {
		if errors.Is(err, models.ErrForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "forbidden",
			})
		}
		if errors.Is(err, models.ErrMessageNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "message not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(msg)
}

// DeleteMessage
// @Summary Видалити повідомлення
// @Tags messages
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/messages/{id} [delete]
func (h *Handler) DeleteMessage(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	messageID := c.Params("id")

	if err := h.services.Message.DeleteMessage(c.Context(), messageID, userID); err != nil {
		if errors.Is(err, models.ErrForbidden) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "forbidden",
			})
		}
		if errors.Is(err, models.ErrMessageNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "message not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "deleted",
	})
}
