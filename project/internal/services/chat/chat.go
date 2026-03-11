package chat

import (
	"context"
	"time"

	"project/internal/models"
)

type chatRepository interface {
	Create(ctx context.Context, chat *models.Chat) (*models.Chat, error)
	FindByID(ctx context.Context, id string) (*models.Chat, error)
	FindPrivateChat(ctx context.Context, id1, id2 string) (*models.Chat, error)
	FindAllByUserID(ctx context.Context, userID string) ([]*models.Chat, error)
	AddMember(ctx context.Context, chatID, userID string) error
	RemoveMember(ctx context.Context, chatID, userID string) error
	UpdateLastMessage(ctx context.Context, chatID string, msg *models.Message) error
	FindByOwnerAndName(ctx context.Context, ownerID, name string) (*models.Chat, error)
}
type Service struct {
	repository chatRepository
}

func NewService(repo chatRepository) *Service {
	return &Service{
		repository: repo,
	}
}

func (s *Service) CreatePrivateChat(ctx context.Context, user1, user2 string) (*models.Chat, error) {
	if chat, err := s.repository.FindPrivateChat(ctx, user1, user2); err == nil {
		return chat, nil
	}
	now := time.Now()
	c, err := s.repository.Create(ctx, &models.Chat{
		Type:      models.ChatTypePrivate,
		Members:   []string{user1, user2},
		OwnerID:   user1, //TODO:
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) CreateGroupChat(ctx context.Context, ownerID, name, description string, members []string) (*models.Chat, error) {

	if chat, err := s.repository.FindByOwnerAndName(ctx, ownerID, name); err == nil {
		return chat, nil
	}
	ownerInMembers := false
	for _, member := range members {
		if member == ownerID {
			ownerInMembers = true
			break
		}
	}
	if !ownerInMembers {
		members = append(members, ownerID)
	}

	now := time.Now()
	c, err := s.repository.Create(ctx, &models.Chat{
		Type:        models.ChatTypeGroup,
		Name:        name,
		Description: description,
		Members:     members,
		OwnerID:     ownerID,
		CreatedAt:   now,
		UpdatedAt:   now,
	})
	if err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetChat(ctx context.Context, chatID, userID string) (*models.Chat, error) {
	chat, err := s.repository.FindByID(ctx, chatID)
	if err != nil {
		return nil, err
	}

	for _, member := range chat.Members {
		if member == userID {
			return chat, nil
		}
	}
	return nil, models.ErrForbidden
}

func (s *Service) GetUserChats(ctx context.Context, userID string) ([]*models.Chat, error) {
	chats, err := s.repository.FindAllByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (s *Service) AddMember(ctx context.Context, chatID, ownerID, userID string) error {
	chat, err := s.GetChat(ctx, chatID, ownerID)
	if err != nil {
		return err
	}

	if chat.Type != models.ChatTypeGroup {
		return models.ErrForbidden
	}

	return s.repository.AddMember(ctx, chatID, userID)
}

func (s *Service) RemoveMember(ctx context.Context, chatID, ownerID, userID string) error {
	chat, err := s.GetChat(ctx, chatID, ownerID)
	if err != nil {
		return err
	}
	if chat.Type != models.ChatTypeGroup {
		return models.ErrForbidden
	}
	if ownerID == userID {
		return models.ErrForbidden
	}
	return s.repository.RemoveMember(ctx, chatID, userID)
}
