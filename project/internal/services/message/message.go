package message

import (
	"context"
	"time"

	"project/internal/models"
)

type messageRepository interface {
	Create(ctx context.Context, msg *models.Message) (*models.Message, error)
	FindByID(ctx context.Context, id string) (*models.Message, error)
	FindByChatID(ctx context.Context, chatID, lastSeeMsgId string, limit int) ([]*models.Message, error)
	UpdateMessage(ctx context.Context, msg *models.Message) error
	SoftDelete(ctx context.Context, msgID string) error
	GetMessage(ctx context.Context, messageID string) (*models.Message, error)
}
type chatRepository interface {
	FindByID(ctx context.Context, chatID string) (*models.Chat, error)
	UpdateLastMessage(ctx context.Context, chatID string, msg *models.Message) error
}

type Service struct {
	chat    chatRepository
	message messageRepository
}

func NewService(message messageRepository, chat chatRepository) *Service {
	return &Service{
		chat:    chat,
		message: message,
	}
}

func checkUserInChat(usrID string, members []string) bool {
	userInChat := false
	for _, member := range members {
		if member == usrID {
			userInChat = true
			break
		}
	}
	return userInChat
}

func (s *Service) Create(ctx context.Context, chatID, userID, text string) (*models.Message, error) {
	chat, err := s.chat.FindByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	userInChat := checkUserInChat(userID, chat.Members)
	if !userInChat {
		return nil, models.ErrForbidden
	}

	now := time.Now()
	msg, err := s.message.Create(ctx, &models.Message{
		ChatID:    chatID,
		UserID:    userID,
		Type:      models.MessageTypeText,
		Text:      text,
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}
	err = s.chat.UpdateLastMessage(ctx, chatID, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *Service) GetMessages(ctx context.Context, chatID, userID, before string, limit int) ([]*models.Message, error) {
	chat, err := s.chat.FindByID(ctx, chatID)
	if err != nil {
		return nil, err
	}
	userInChat := checkUserInChat(userID, chat.Members)
	if !userInChat {
		return nil, models.ErrForbidden
	}
	msgs, err := s.message.FindByChatID(ctx, chatID, before, limit)
	if err != nil {
		return nil, err
	}
	return msgs, nil
}

func (s *Service) EditMessage(ctx context.Context, messageID, userID, text string) (*models.Message, error) {
	msg, err := s.message.FindByID(ctx, messageID)
	if err != nil {
		return nil, err
	}

	if msg.UserID != userID {
		return nil, models.ErrForbidden
	}

	if msg.IsDeleted {
		return nil, models.ErrForbidden
	}
	now := time.Now()
	msg.Text = text
	msg.IsEdited = true
	msg.UpdatedAt = now
	err = s.message.UpdateMessage(ctx, msg)
	if err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *Service) DeleteMessage(ctx context.Context, messageID, userID string) error {
	msg, err := s.message.FindByID(ctx, messageID)
	if err != nil {
		return err
	}

	if msg.UserID != userID {
		return models.ErrForbidden
	}

	if msg.IsDeleted {
		return models.ErrForbidden
	}

	return s.message.SoftDelete(ctx, messageID)
}

func (s *Service) GetMessage(ctx context.Context, messageID string) (*models.Message, error) {
	msg, err := s.message.FindByID(ctx, messageID)
	if err != nil {
		return nil, err
	}
	return msg, nil
}
