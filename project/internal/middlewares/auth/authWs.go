package auth

import (
	"strings"

	"github.com/gofiber/fiber/v3"
)

func (m *Middleware) HandleWs(c fiber.Ctx) error {
	if allowed := c.Locals("allowed"); allowed == nil {
	}

	token := strings.Trim(c.Query("token"), "\"")
	if token == "" {
		return fiber.ErrUnauthorized
	}

	claims, err := m.service.ParseAuthToken(token)
	if err != nil {
		return fiber.ErrUnauthorized
	}
	c.Locals("user_id", claims.UserID)
	c.Locals("token_exp", claims.ExpiresAt)
	return c.Next()
}
