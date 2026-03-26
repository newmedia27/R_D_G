package handlers

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/v3/swaggerui"
	"github.com/gofiber/contrib/v3/websocket"
	"github.com/gofiber/fiber/v3"
	"project/internal/config"
	"project/internal/handlers/auth"
	"project/internal/handlers/call"
	"project/internal/handlers/chat"
	"project/internal/handlers/message"
	"project/internal/handlers/user"
	"project/internal/middlewares"
	"project/internal/services"
	"project/internal/ws"
)

type Handlers struct {
	middlewares *middlewares.Middlewares
	hub         *ws.Hub

	Auth    *auth.Handler
	Users   *user.Handler
	Chat    *chat.Handler
	Message *message.Handler
	Ws      *ws.Handler
	Call    *call.Handler
}

func NewHandlers(m *middlewares.Middlewares, s *services.Services, cfg *config.Config, hub *ws.Hub) *Handlers {
	v := validator.New()
	return &Handlers{
		middlewares: m,
		hub:         hub,

		Auth:    auth.NewHandler(s, v, cfg),
		Users:   user.NewHandler(s, v),
		Chat:    chat.NewHandler(s, v, hub),
		Message: message.NewHandler(s, v),
		Ws:      ws.NewHandler(s, hub),
		Call:    call.NewHandler(s, cfg),
	}
}

// Health check
// @Summary Health check
// @Description Перевіряє стан сервера та БД
// @Tags health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /health [get]
func (h *Handlers) Health(ctx fiber.Ctx) error {
	return ctx.JSON(fiber.Map{
		"status":    "ok",
		"version":   "v1",
		"timestamp": time.Now().UTC().Format(time.RFC3339Nano),
	})
}

func (h *Handlers) InitRouters(app *fiber.App) {
	app.Use(h.middlewares.Recover)
	app.Use(h.middlewares.Cors)
	app.Get("/health", h.Health)

	app.Use(swaggerui.New(swaggerui.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
	}))

	api := app.Group("/api/v1")
	api.Use(h.middlewares.Log)
	//WS
	api.Use("/ws", func(c fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	api.Get("/ws", h.middlewares.Auth.HandleWs, websocket.New(h.Ws.Handle))

	Auth := api.Group("/auth")
	{
		Auth.Post("/signin", h.Auth.Login)
		Auth.Post("/signup", h.Auth.Register)
		Auth.Post("/refresh", h.Auth.Refresh)
		Auth.Delete("/logout", h.Auth.Logout)
	}

	protected := api.Group("/")
	{
		protected.Use(h.middlewares.Auth.Handle)
		protected.Delete("/auth/logout/all", h.Auth.LogoutAll)

		protected.Get("/users", h.Users.GetUsers)
		protected.Get("/users/search", h.Users.Search)
		protected.Post("/users/batch", h.Users.GetUsersByIds)
		protected.Get("/users/profile", h.Users.GetProfile)
		protected.Get("/users/:id", h.Users.GetUser)

		// Chats
		protected.Post("/chats/group", h.Chat.CreateGroup)
		protected.Post("/chats/private/:userId", h.Chat.CreatePrivate)
		protected.Get("/chats", h.Chat.GetUserChats)
		protected.Get("/chats/:id", h.Chat.GetChat)
		protected.Post("/chats/:id/members", h.Chat.AddMember)
		protected.Delete("/chats/:id/members/:userId", h.Chat.RemoveMember)

		//call
		protected.Post("/calls/token", h.Call.GetToken)

		// Messages
		protected.Get("/chats/:id/messages", h.Message.GetMessages)

	}

	api.Get("/", func(ctx fiber.Ctx) {
		_ = ctx.SendString("Hello, World!")
	})

	//app.Get("/swagger/*", swaggerui.New(swaggerui.Config{
	//	BasePath: "/swagger",
	//	FilePath: "./docs/swagger.json",
	//}))

}
