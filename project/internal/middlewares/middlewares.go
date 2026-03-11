package middlewares

import (
	"github.com/gofiber/fiber/v3"
	"project/internal/config"
	"project/internal/middlewares/auth"
	"project/internal/middlewares/cors"
	"project/internal/middlewares/log"
	rec "project/internal/middlewares/recover"
	"project/internal/services"
)

type Middlewares struct {
	Log     fiber.Handler
	Cors    fiber.Handler
	Recover fiber.Handler
	Auth    *auth.Middleware
}

func NewMiddlewares(cfg *config.Config, s *services.Services) *Middlewares {
	return &Middlewares{
		Log:     log.Middleware(),
		Cors:    cors.Middleware(cfg),
		Recover: rec.Middleware(),
		Auth:    auth.NewMiddleware(s.Auth),
	}
}
