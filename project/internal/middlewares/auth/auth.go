package auth

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"project/internal/services/auth"
)

type authService interface {
	ParseAuthToken(token string) (*auth.TokenClaims, error)
}
type Middleware struct {
	service authService
}

func NewMiddleware(service authService) *Middleware {
	return &Middleware{
		service: service,
	}
}

func (m *Middleware) Handle(ctx fiber.Ctx) error {
	accessToken := ctx.Get("Authorization")
	if accessToken == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error":   "Unauthorized",
				"message": "Invalid token format",
			},
		)
	}
	if !strings.HasPrefix(accessToken, "Bearer ") {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid token format",
		})
	}
	token := accessToken[7:]
	if token == "" {
		return ctx.Status(fiber.StatusUnauthorized).JSON(
			fiber.Map{
				"error":   "Unauthorized",
				"message": "Invalid token",
			},
		)
	}

	claims, err := m.service.ParseAuthToken(token)
	if err != nil {
		return fiber.ErrUnauthorized
	}

	ctx.Locals("user_id", claims.UserID)

	return ctx.Next()
}
