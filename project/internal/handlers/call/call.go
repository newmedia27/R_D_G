package call

import (
	"github.com/gofiber/fiber/v3"
	"project/internal/config"
	"project/internal/services"
)

type Handler struct {
	services *services.Services
	cfg      *config.Config
}

func NewHandler(s *services.Services, cfg *config.Config) *Handler {
	return &Handler{
		services: s,
		cfg:      cfg,
	}
}

type GetTokenRequest struct {
	RoomName string `json:"chat_id"`
}

func (h *Handler) GetToken(c fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	var req GetTokenRequest

	err := c.Bind().JSON(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid chat_id",
		})
	}
	token, err := h.services.Call.GetToken(userID, req.RoomName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
		"url":   h.cfg.LiveKitURL,
	})
}
