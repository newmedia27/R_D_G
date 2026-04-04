package user

import (
	"context"

	"project/internal/models"
)

type userRepository interface {
	FindByIds(ctx context.Context, ids []string) ([]*models.User, error)
	Search(ctx context.Context, userId, query string) ([]*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
}
type Service struct {
	user userRepository
}

func NewService(user userRepository) *Service {
	return &Service{
		user: user,
	}
}

func (s *Service) GetByIds(ctx context.Context, ids []string) (map[string]*models.User, error) {
	users, err := s.user.FindByIds(ctx, ids)
	if err != nil {
		return nil, err
	}
	result := make(map[string]*models.User, len(users))
	for _, user := range users {
		result[user.ID] = user
	}
	return result, nil
}

func (s *Service) Search(ctx context.Context, userID, query string) ([]*models.User, error) {
	return s.user.Search(ctx, userID, query)
}

func (s *Service) FindByID(ctx context.Context, id string) (*models.User, error) {
	return s.user.FindByID(ctx, id)
}
