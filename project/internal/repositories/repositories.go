package repositories

import (
	"project/internal/clients"
	"project/internal/repositories/call"
	"project/internal/repositories/chat"
	"project/internal/repositories/message"
	"project/internal/repositories/session"
	"project/internal/repositories/user"
)

type Repositories struct {
	Session *session.Repository
	User    *user.Repository
	Chat    *chat.Repository
	Message *message.Repository
	Call    *call.Repository
}

func NewRepositories(clients *clients.Clients) *Repositories {
	mongo := clients.Mongo.Db
	lks := clients.LiveKit
	return &Repositories{
		Session: session.NewRepository(mongo),
		User:    user.NewRepository(mongo),
		Chat:    chat.NewRepository(mongo),
		Message: message.NewRepository(mongo),
		Call:    call.NewRepository(lks),
	}
}
