package services

import (
	"project/internal/config"
	"project/internal/repositories"
	"project/internal/services/auth"
	"project/internal/services/call"
	"project/internal/services/chat"
	"project/internal/services/message"
	"project/internal/services/user"
)

type Services struct {
	Auth    *auth.Service
	Chat    *chat.Service
	Message *message.Service
	User    *user.Service
	Call    *call.Service
}

func NewServices(cfg *config.Config, repo *repositories.Repositories) *Services {
	return &Services{
		Auth:    auth.NewService(repo.User, repo.Session, cfg.JWTSecret, cfg.JWTExpiration, cfg.RefreshExpiration),
		Chat:    chat.NewService(repo.Chat),
		Message: message.NewService(repo.Message, repo.Chat),
		User:    user.NewService(repo.User),
		Call:    call.NewService(repo.Call),
	}
}
